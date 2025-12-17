import { stats } from '@nemo-network-indexer/base/build';
import {
  SubaccountTable,
  SubaccountFromDatabase,
  FillTable,
  FillFromDatabase,
  FillColumns,
  Liquidity,
  BlockTable,
  BlockFromDatabase,
  FundingIndexUpdatesTable,
  FundingIndexMap,
  PerpetualPositionTable,
  PerpetualPositionFromDatabase,
  PerpetualPositionStatus,
  AssetPositionTable,
  AssetPositionFromDatabase,
  AssetTable,
  AssetFromDatabase,
  MarketTable,
  MarketFromDatabase,
  perpetualMarketRefresher,
  helpers,
  QueryableField,
} from '@nemo-network-indexer/postgres/build/src';
import Big from 'big.js';
import express from 'express';
import { matchedData } from 'express-validator';
import { DateTime } from 'luxon';
import {
  Controller, Get, Path, Route,
} from 'tsoa';

import { getReqRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import { complianceAndGeoCheck } from '../../../lib/compliance-and-geo-check';
import { NotFoundError } from '../../../lib/errors';
import {
  getFundingIndexMaps,
  handleControllerError,
  getSubaccountResponse,
  getTotalUnsettledFunding,
} from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import { CheckAddressSchema } from '../../../lib/validation/schemas';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import {
  PortfolioValueResponse,
  VolumeResponse,
  FeesPercentageResponse,
  TotalFundingFeeResponse,
} from '../../../types';

const router: express.Router = express.Router();
const controllerName: string = 'portfolio-controller';

@Route('portfolio')
class PortfolioController extends Controller {
  @Get('/:address/value')
  public async getPortfolioValue(
    @Path() address: string,
  ): Promise<PortfolioValueResponse> {
    const subaccounts: SubaccountFromDatabase[] = await SubaccountTable.findAll(
      {
        address,
      },
      [],
    );

    if (subaccounts.length === 0) {
      throw new NotFoundError(`No subaccounts found for address ${address}`);
    }

    const [latestBlock, assets, markets] = await Promise.all([
      BlockTable.getLatest(),
      AssetTable.findAll({}, []),
      MarketTable.findAll({}, []),
    ]);

    const latestFundingIndexMap: FundingIndexMap = await FundingIndexUpdatesTable
      .findFundingIndexMap(latestBlock.blockHeight);

    const perpetualMarketsMap = perpetualMarketRefresher.getPerpetualMarketsMap();

    let totalEquity = Big(0);

    for (const subaccount of subaccounts) {
      const [
        perpetualPositions,
        assetPositions,
        lastUpdatedFundingIndexMap,
      ] = await Promise.all([
        PerpetualPositionTable.findAll(
          {
            subaccountId: [subaccount.id],
            status: [PerpetualPositionStatus.OPEN],
          },
          [],
        ),
        AssetPositionTable.findAll(
          {
            subaccountId: [subaccount.id],
          },
          [],
        ),
        FundingIndexUpdatesTable.findFundingIndexMap(subaccount.updatedAtHeight),
      ]);

      const subaccountResponse = getSubaccountResponse(
        subaccount,
        perpetualPositions,
        assetPositions,
        assets,
        markets,
        perpetualMarketsMap,
        latestBlock.blockHeight,
        latestFundingIndexMap,
        lastUpdatedFundingIndexMap,
      );

      totalEquity = totalEquity.plus(subaccountResponse.equity);
    }

    return {
      address,
      portfolioValue: totalEquity.toFixed(),
    };
  }

  @Get('/:address/volume')
  public async getVolume(
    @Path() address: string,
  ): Promise<VolumeResponse> {
    const subaccounts: SubaccountFromDatabase[] = await SubaccountTable.findAll(
      {
        address,
      },
      [],
    );

    if (subaccounts.length === 0) {
      throw new NotFoundError(`No subaccounts found for address ${address}`);
    }

    const subaccountIds = subaccounts.map((s) => s.id);
    const thirtyDaysAgo = DateTime.now().minus({ days: 30 }).toISO();

    const { results: fills } = await FillTable.findAll(
      {
        subaccountId: subaccountIds,
        createdOnOrAfter: thirtyDaysAgo,
      },
      [QueryableField.CREATED_ON_OR_AFTER],
    );

    let totalVolume = Big(0);
    for (const fill of fills) {
      // Use quoteAmount if available, otherwise calculate as price * size
      const volume = fill.quoteAmount
        ? Big(fill.quoteAmount)
        : Big(fill.price).times(fill.size);
      totalVolume = totalVolume.plus(volume);
    }

    return {
      address,
      volume: totalVolume.toFixed(),
      periodDays: 30,
    };
  }

  @Get('/:address/fees-percentage')
  public async getFeesPercentage(
    @Path() address: string,
  ): Promise<FeesPercentageResponse> {
    const subaccounts: SubaccountFromDatabase[] = await SubaccountTable.findAll(
      {
        address,
      },
      [],
    );

    if (subaccounts.length === 0) {
      throw new NotFoundError(`No subaccounts found for address ${address}`);
    }

    const subaccountIds = subaccounts.map((s) => s.id);

    const { results: fills } = await FillTable.findAll(
      {
        subaccountId: subaccountIds,
      },
      [],
    );

    let takerFees = Big(0);
    let makerFees = Big(0);
    let totalFees = Big(0);

    for (const fill of fills) {
      const fee = Big(fill.fee || '0');
      totalFees = totalFees.plus(fee);

      if (fill.liquidity === Liquidity.TAKER) {
        takerFees = takerFees.plus(fee);
      } else if (fill.liquidity === Liquidity.MAKER) {
        makerFees = makerFees.plus(fee);
      }
    }

    const takerPercentage = totalFees.gt(0)
      ? takerFees.div(totalFees).times(100)
      : Big(0);
    const makerPercentage = totalFees.gt(0)
      ? makerFees.div(totalFees).times(100)
      : Big(0);

    return {
      address,
      takerPercentage: takerPercentage.toFixed(),
      makerPercentage: makerPercentage.toFixed(),
      totalFees: totalFees.toFixed(),
    };
  }

  @Get('/:address/total-funding-fee')
  public async getTotalFundingFee(
    @Path() address: string,
  ): Promise<TotalFundingFeeResponse> {
    const subaccounts: SubaccountFromDatabase[] = await SubaccountTable.findAll(
      {
        address,
      },
      [],
    );

    if (subaccounts.length === 0) {
      throw new NotFoundError(`No subaccounts found for address ${address}`);
    }

    const [latestBlock, assets, markets] = await Promise.all([
      BlockTable.getLatest(),
      AssetTable.findAll({}, []),
      MarketTable.findAll({}, []),
    ]);

    const latestFundingIndexMap: FundingIndexMap = await FundingIndexUpdatesTable
      .findFundingIndexMap(latestBlock.blockHeight);

    const perpetualMarketsMap = perpetualMarketRefresher.getPerpetualMarketsMap();

    let totalUnsettledFunding = Big(0);
    let totalRealizedFunding = Big(0);

    for (const subaccount of subaccounts) {
      const [
        perpetualPositions,
        assetPositions,
        lastUpdatedFundingIndexMap,
      ] = await Promise.all([
        PerpetualPositionTable.findAll(
          {
            subaccountId: [subaccount.id],
            status: [PerpetualPositionStatus.OPEN],
          },
          [],
        ),
        AssetPositionTable.findAll(
          {
            subaccountId: [subaccount.id],
          },
          [],
        ),
        FundingIndexUpdatesTable.findFundingIndexMap(subaccount.updatedAtHeight),
      ]);

      // Calculate unsettled funding for open positions
      const unsettledFunding = getTotalUnsettledFunding(
        perpetualPositions,
        latestFundingIndexMap,
        lastUpdatedFundingIndexMap,
      );
      totalUnsettledFunding = totalUnsettledFunding.plus(unsettledFunding);

    }

    // Total funding fee = unsettled funding (negative means paid, positive means received)
    // We'll return it as a positive value if paid, negative if received
    const totalFundingFee = totalUnsettledFunding.neg();

    return {
      address,
      totalFundingFee: totalFundingFee.toFixed(),
      unsettledFunding: totalUnsettledFunding.toFixed(),
    };
  }
}

router.get(
  '/:address/value',
  rateLimiterMiddleware(getReqRateLimiter),
  ...CheckAddressSchema,
  handleValidationErrors,
  complianceAndGeoCheck,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      address,
    }: {
      address: string,
    } = matchedData(req) as { address: string };

    try {
      const controller: PortfolioController = new PortfolioController();
      const response: PortfolioValueResponse = await controller.getPortfolioValue(
        address,
      );

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'PortfolioController GET /:address/value',
        'Portfolio value error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_portfolio_value.timing`,
        Date.now() - start,
      );
    }
  },
);

router.get(
  '/:address/volume',
  rateLimiterMiddleware(getReqRateLimiter),
  ...CheckAddressSchema,
  handleValidationErrors,
  complianceAndGeoCheck,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      address,
    }: {
      address: string,
    } = matchedData(req) as { address: string };

    try {
      const controller: PortfolioController = new PortfolioController();
      const response: VolumeResponse = await controller.getVolume(
        address,
      );

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'PortfolioController GET /:address/volume',
        'Volume error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_volume.timing`,
        Date.now() - start,
      );
    }
  },
);

router.get(
  '/:address/fees-percentage',
  rateLimiterMiddleware(getReqRateLimiter),
  ...CheckAddressSchema,
  handleValidationErrors,
  complianceAndGeoCheck,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      address,
    }: {
      address: string,
    } = matchedData(req) as { address: string };

    try {
      const controller: PortfolioController = new PortfolioController();
      const response: FeesPercentageResponse = await controller.getFeesPercentage(
        address,
      );

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'PortfolioController GET /:address/fees-percentage',
        'Fees percentage error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_fees_percentage.timing`,
        Date.now() - start,
      );
    }
  },
);

router.get(
  '/:address/total-funding-fee',
  rateLimiterMiddleware(getReqRateLimiter),
  ...CheckAddressSchema,
  handleValidationErrors,
  complianceAndGeoCheck,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      address,
    }: {
      address: string,
    } = matchedData(req) as { address: string };

    try {
      const controller: PortfolioController = new PortfolioController();
      const response: TotalFundingFeeResponse = await controller.getTotalFundingFee(
        address,
      );

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'PortfolioController GET /:address/total-funding-fee',
        'Total funding fee error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_total_funding_fee.timing`,
        Date.now() - start,
      );
    }
  },
);

export default router;

