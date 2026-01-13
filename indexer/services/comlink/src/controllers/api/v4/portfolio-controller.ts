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
  PnlTicksTable,
  PnlTicksFromDatabase,
  IsoString,
  Ordering,
  DEFAULT_POSTGRES_OPTIONS,
  OrderSide,
  PositionSide,
} from '@nemo-network-indexer/postgres/build/src';
import Big from 'big.js';
import express from 'express';
import { matchedData } from 'express-validator';
import { DateTime } from 'luxon';
import {
  Controller, Get, Path, Query, Route,
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
import { CheckAddressSchema, CheckSubaccountSchema, CheckLimitAndCreatedBeforeOrAtAndOnOrAfterSchema } from '../../../lib/validation/schemas';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import {
  PortfolioValueResponse,
  VolumeResponse,
  FeesPercentageResponse,
  TotalFundingFeeResponse,
  LivePnlResponse,
  RealizedPnlResponse,
  ProfitFactorResponse,
  MaxDrawdownResponse,
  HealthResponse,
  EquityResponse,
  EquityListResponse,
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
    let totalFeesPaid = Big(0); // Sum of absolute fee values for percentage calculation

    for (const fill of fills) {
      const fee = Big(fill.fee || '0');
      const absFee = fee.abs();

      if (fill.liquidity === Liquidity.TAKER) {
        takerFees = takerFees.plus(fee);
        totalFeesPaid = totalFeesPaid.plus(absFee);
      } else if (fill.liquidity === Liquidity.MAKER) {
        makerFees = makerFees.plus(fee);
        totalFeesPaid = totalFeesPaid.plus(absFee);
      }
    }

    // Calculate percentages based on absolute fee amounts
    // This ensures percentages add up to 100% correctly
    const takerAbsFees = takerFees.abs();
    const makerAbsFees = makerFees.abs();
    const totalAbsFees = takerAbsFees.plus(makerAbsFees);

    const takerPercentage = totalAbsFees.gt(0)
      ? takerAbsFees.div(totalAbsFees).times(100)
      : Big(0);
    const makerPercentage = totalAbsFees.gt(0)
      ? makerAbsFees.div(totalAbsFees).times(100)
      : Big(0);

    // Total fees is the net (taker fees paid - maker rebates received)
    const totalFees = takerFees.plus(makerFees);

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

  @Get('/:address/:subaccountNumber/live-pnl')
  public async getLivePnl(
    @Path() address: string,
    @Path() subaccountNumber: number,
    @Query() createdOnOrAfter?: IsoString,
    @Query() createdBeforeOrAt?: IsoString,
  ): Promise<LivePnlResponse> {
    const subaccountId: string = SubaccountTable.uuid(address, subaccountNumber);
    const subaccount: SubaccountFromDatabase | undefined = await SubaccountTable.findById(subaccountId);

    if (subaccount === undefined) {
      throw new NotFoundError(`No subaccount found with address ${address} and subaccountNumber ${subaccountNumber}`);
    }

    // Get latest PnL tick within date range or calculate from current positions
    const requiredFields: QueryableField[] = [];
    const queryConfig: {
      subaccountId: string[];
      createdOnOrAfter?: IsoString;
      createdBeforeOrAt?: IsoString;
      limit: number;
    } = {
      subaccountId: [subaccountId],
      limit: 1,
    };

    if (createdOnOrAfter) {
      requiredFields.push(QueryableField.CREATED_ON_OR_AFTER);
      queryConfig.createdOnOrAfter = createdOnOrAfter;
    }
    if (createdBeforeOrAt) {
      requiredFields.push(QueryableField.CREATED_BEFORE_OR_AT);
      queryConfig.createdBeforeOrAt = createdBeforeOrAt;
    }

    const { results: pnlTicks } = await PnlTicksTable.findAll(
      queryConfig,
      requiredFields.length > 0 ? requiredFields : [QueryableField.LIMIT],
      {
        ...DEFAULT_POSTGRES_OPTIONS,
        orderBy: [[QueryableField.BLOCK_HEIGHT, Ordering.DESC]],
      },
    );
    const latestPnlTick: PnlTicksFromDatabase | undefined = pnlTicks[0];
    
    let livePnl = Big(0);
    if (latestPnlTick) {
      livePnl = Big(latestPnlTick.totalPnl);
    } else {
      // Calculate from current positions if no PnL tick exists
      const [latestBlock, assets, markets] = await Promise.all([
        BlockTable.getLatest(),
        AssetTable.findAll({}, []),
        MarketTable.findAll({}, []),
      ]);

      const latestFundingIndexMap: FundingIndexMap = await FundingIndexUpdatesTable
        .findFundingIndexMap(latestBlock.blockHeight);

      const perpetualMarketsMap = perpetualMarketRefresher.getPerpetualMarketsMap();

      const [
        perpetualPositions,
        assetPositions,
        lastUpdatedFundingIndexMap,
      ] = await Promise.all([
        PerpetualPositionTable.findAll(
          {
            subaccountId: [subaccountId],
            status: [PerpetualPositionStatus.OPEN],
          },
          [],
        ),
        AssetPositionTable.findAll(
          {
            subaccountId: [subaccountId],
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

      // Calculate PnL as equity minus net transfers
      // This is a simplified calculation - actual PnL would need historical data
      livePnl = Big(subaccountResponse.equity);
    }

    return {
      address,
      subaccountNumber,
      value: livePnl.toFixed(),
    };
  }

  @Get('/:address/:subaccountNumber/realized-pnl')
  public async getRealizedPnl(
    @Path() address: string,
    @Path() subaccountNumber: number,
    @Query() createdOnOrAfter?: IsoString,
    @Query() createdBeforeOrAt?: IsoString,
  ): Promise<RealizedPnlResponse> {
    const subaccountId: string = SubaccountTable.uuid(address, subaccountNumber);
    const subaccount: SubaccountFromDatabase | undefined = await SubaccountTable.findById(subaccountId);

    if (subaccount === undefined) {
      throw new NotFoundError(`No subaccount found with address ${address} and subaccountNumber ${subaccountNumber}`);
    }

    // Get fills in date range
    const requiredFields: QueryableField[] = [];
    const queryConfig: {
      subaccountId: string[];
      createdOnOrAfter?: IsoString;
      createdBeforeOrAt?: IsoString;
    } = {
      subaccountId: [subaccountId],
    };

    if (createdOnOrAfter) {
      requiredFields.push(QueryableField.CREATED_ON_OR_AFTER);
      queryConfig.createdOnOrAfter = createdOnOrAfter;
    }
    if (createdBeforeOrAt) {
      requiredFields.push(QueryableField.CREATED_BEFORE_OR_AT);
      queryConfig.createdBeforeOrAt = createdBeforeOrAt;
    }

    const { results: fills } = await FillTable.findAll(
      queryConfig,
      requiredFields,
    );

    // Get all perpetual positions to track entry prices
    const perpetualPositions: PerpetualPositionFromDatabase[] = await PerpetualPositionTable.findAll(
      {
        subaccountId: [subaccountId],
      },
      [],
    );

    const perpetualMarketsMap = perpetualMarketRefresher.getPerpetualMarketsMap();
    const positionMap: Record<string, PerpetualPositionFromDatabase> = {};
    perpetualPositions.forEach((position: PerpetualPositionFromDatabase) => {
      const perpetualMarket = perpetualMarketRefresher.getPerpetualMarketFromId(position.perpetualId);
      if (perpetualMarket) {
        const key = `${position.subaccountId}-${perpetualMarket.clobPairId}`;
        positionMap[key] = position;
      }
    });

    let realizedPnl = Big(0);

    for (const fill of fills) {
      const perpetualMarket = perpetualMarketRefresher.getPerpetualMarketFromClobPairId(fill.clobPairId);
      if (!perpetualMarket) continue;

      const key = `${fill.subaccountId}-${fill.clobPairId}`;
      const position = positionMap[key];

      if (!position) continue;

      // Check if fill closes the position
      const isClosingLong = position.side === PositionSide.LONG && fill.side === OrderSide.SELL;
      const isClosingShort = position.side === PositionSide.SHORT && fill.side === OrderSide.BUY;

      if (isClosingLong || isClosingShort) {
        const fillPrice = Big(fill.price);
        const fillSize = Big(fill.size);
        const entryPrice = Big(position.entryPrice);
        const fee = Big(fill.fee || '0');

        let closedPnL: Big;
        if (isClosingLong) {
          closedPnL = fillPrice.minus(entryPrice).mul(fillSize).minus(fee);
        } else {
          closedPnL = entryPrice.minus(fillPrice).mul(fillSize).minus(fee);
        }

        realizedPnl = realizedPnl.plus(closedPnL);
      }
    }

    return {
      address,
      subaccountNumber,
      value: realizedPnl.toFixed(),
    };
  }

  @Get('/:address/:subaccountNumber/profit-factor')
  public async getProfitFactor(
    @Path() address: string,
    @Path() subaccountNumber: number,
    @Query() createdOnOrAfter?: IsoString,
    @Query() createdBeforeOrAt?: IsoString,
  ): Promise<ProfitFactorResponse> {
    const subaccountId: string = SubaccountTable.uuid(address, subaccountNumber);
    const subaccount: SubaccountFromDatabase | undefined = await SubaccountTable.findById(subaccountId);

    if (subaccount === undefined) {
      throw new NotFoundError(`No subaccount found with address ${address} and subaccountNumber ${subaccountNumber}`);
    }

    // Get fills in date range
    const requiredFields: QueryableField[] = [];
    const queryConfig: {
      subaccountId: string[];
      createdOnOrAfter?: IsoString;
      createdBeforeOrAt?: IsoString;
    } = {
      subaccountId: [subaccountId],
    };

    if (createdOnOrAfter) {
      requiredFields.push(QueryableField.CREATED_ON_OR_AFTER);
      queryConfig.createdOnOrAfter = createdOnOrAfter;
    }
    if (createdBeforeOrAt) {
      requiredFields.push(QueryableField.CREATED_BEFORE_OR_AT);
      queryConfig.createdBeforeOrAt = createdBeforeOrAt;
    }

    const { results: fills } = await FillTable.findAll(
      queryConfig,
      requiredFields,
    );

    // Get all perpetual positions to track entry prices
    const perpetualPositions: PerpetualPositionFromDatabase[] = await PerpetualPositionTable.findAll(
      {
        subaccountId: [subaccountId],
      },
      [],
    );

    const perpetualMarketsMap = perpetualMarketRefresher.getPerpetualMarketsMap();
    const positionMap: Record<string, PerpetualPositionFromDatabase> = {};
    perpetualPositions.forEach((position: PerpetualPositionFromDatabase) => {
      const perpetualMarket = perpetualMarketRefresher.getPerpetualMarketFromId(position.perpetualId);
      if (perpetualMarket) {
        const key = `${position.subaccountId}-${perpetualMarket.clobPairId}`;
        positionMap[key] = position;
      }
    });

    let grossProfit = Big(0);
    let grossLoss = Big(0);

    for (const fill of fills) {
      const perpetualMarket = perpetualMarketRefresher.getPerpetualMarketFromClobPairId(fill.clobPairId);
      if (!perpetualMarket) continue;

      const key = `${fill.subaccountId}-${fill.clobPairId}`;
      const position = positionMap[key];

      if (!position) continue;

      // Check if fill closes the position
      const isClosingLong = position.side === PositionSide.LONG && fill.side === OrderSide.SELL;
      const isClosingShort = position.side === PositionSide.SHORT && fill.side === OrderSide.BUY;

      if (isClosingLong || isClosingShort) {
        const fillPrice = Big(fill.price);
        const fillSize = Big(fill.size);
        const entryPrice = Big(position.entryPrice);
        const fee = Big(fill.fee || '0');

        let closedPnL: Big;
        if (isClosingLong) {
          closedPnL = fillPrice.minus(entryPrice).mul(fillSize).minus(fee);
        } else {
          closedPnL = entryPrice.minus(fillPrice).mul(fillSize).minus(fee);
        }

        if (closedPnL.gt(0)) {
          grossProfit = grossProfit.plus(closedPnL);
        } else {
          grossLoss = grossLoss.plus(closedPnL.abs());
        }
      }
    }

    const profitFactor = grossLoss.gt(0) ? grossProfit.div(grossLoss) : (grossProfit.gt(0) ? Big(Infinity) : Big(0));

    return {
      address,
      subaccountNumber,
      value: profitFactor.toFixed(),
    };
  }

  @Get('/:address/:subaccountNumber/max-drawdown')
  public async getMaxDrawdown(
    @Path() address: string,
    @Path() subaccountNumber: number,
    @Query() createdOnOrAfter?: IsoString,
    @Query() createdBeforeOrAt?: IsoString,
  ): Promise<MaxDrawdownResponse> {
    const subaccountId: string = SubaccountTable.uuid(address, subaccountNumber);
    const subaccount: SubaccountFromDatabase | undefined = await SubaccountTable.findById(subaccountId);

    if (subaccount === undefined) {
      throw new NotFoundError(`No subaccount found with address ${address} and subaccountNumber ${subaccountNumber}`);
    }

    // Get PnL ticks in date range
    const requiredFields: QueryableField[] = [];
    const queryConfig: {
      subaccountId: string[];
      createdOnOrAfter?: IsoString;
      createdBeforeOrAt?: IsoString;
    } = {
      subaccountId: [subaccountId],
    };

    if (createdOnOrAfter) {
      requiredFields.push(QueryableField.CREATED_ON_OR_AFTER);
      queryConfig.createdOnOrAfter = createdOnOrAfter;
    }
    if (createdBeforeOrAt) {
      requiredFields.push(QueryableField.CREATED_BEFORE_OR_AT);
      queryConfig.createdBeforeOrAt = createdBeforeOrAt;
    }

    const { results: pnlTicks } = await PnlTicksTable.findAll(
      queryConfig,
      requiredFields,
      {
        ...DEFAULT_POSTGRES_OPTIONS,
        orderBy: [[QueryableField.BLOCK_HEIGHT, Ordering.ASC]],
      },
    );

    if (pnlTicks.length === 0) {
      return {
        address,
        subaccountNumber,
        value: '0',
      };
    }

    let maxDrawdown = Big(0);
    let peakEquity = Big(pnlTicks[0].equity);

    for (const tick of pnlTicks) {
      const equity = Big(tick.equity);
      
      if (equity.gt(peakEquity)) {
        peakEquity = equity;
      }

      const drawdown = peakEquity.minus(equity);
      if (drawdown.gt(maxDrawdown)) {
        maxDrawdown = drawdown;
      }
    }

    return {
      address,
      subaccountNumber,
      value: maxDrawdown.toFixed(),
    };
  }

  @Get('/:address/:subaccountNumber/health')
  public async getHealth(
    @Path() address: string,
    @Path() subaccountNumber: number,
  ): Promise<HealthResponse> {
    const subaccountId: string = SubaccountTable.uuid(address, subaccountNumber);
    const subaccount: SubaccountFromDatabase | undefined = await SubaccountTable.findById(subaccountId);

    if (subaccount === undefined) {
      throw new NotFoundError(`No subaccount found with address ${address} and subaccountNumber ${subaccountNumber}`);
    }

    const [latestBlock, assets, markets] = await Promise.all([
      BlockTable.getLatest(),
      AssetTable.findAll({}, []),
      MarketTable.findAll({}, []),
    ]);

    const latestFundingIndexMap: FundingIndexMap = await FundingIndexUpdatesTable
      .findFundingIndexMap(latestBlock.blockHeight);

    const perpetualMarketsMap = perpetualMarketRefresher.getPerpetualMarketsMap();

    const [
      perpetualPositions,
      assetPositions,
      lastUpdatedFundingIndexMap,
    ] = await Promise.all([
      PerpetualPositionTable.findAll(
        {
          subaccountId: [subaccountId],
          status: [PerpetualPositionStatus.OPEN],
        },
        [],
      ),
      AssetPositionTable.findAll(
        {
          subaccountId: [subaccountId],
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

    // Health is typically calculated as equity / margin requirements
    // Using freeCollateral / equity as a health metric (higher is better)
    const equity = Big(subaccountResponse.equity);
    const freeCollateral = Big(subaccountResponse.freeCollateral);
    
    const health = equity.gt(0) ? freeCollateral.div(equity).times(100) : Big(0);

    return {
      address,
      subaccountNumber,
      value: health.toFixed(),
    };
  }

  @Get('/:address/:subaccountNumber/equity')
  public async getEquity(
    @Path() address: string,
    @Path() subaccountNumber: number,
  ): Promise<EquityResponse> {
    const subaccountId: string = SubaccountTable.uuid(address, subaccountNumber);
    const subaccount: SubaccountFromDatabase | undefined = await SubaccountTable.findById(subaccountId);

    if (subaccount === undefined) {
      throw new NotFoundError(`No subaccount found with address ${address} and subaccountNumber ${subaccountNumber}`);
    }

    const [latestBlock, assets, markets] = await Promise.all([
      BlockTable.getLatest(),
      AssetTable.findAll({}, []),
      MarketTable.findAll({}, []),
    ]);

    const latestFundingIndexMap: FundingIndexMap = await FundingIndexUpdatesTable
      .findFundingIndexMap(latestBlock.blockHeight);

    const perpetualMarketsMap = perpetualMarketRefresher.getPerpetualMarketsMap();

    const [
      perpetualPositions,
      assetPositions,
      lastUpdatedFundingIndexMap,
    ] = await Promise.all([
      PerpetualPositionTable.findAll(
        {
          subaccountId: [subaccountId],
          status: [PerpetualPositionStatus.OPEN],
        },
        [],
      ),
      AssetPositionTable.findAll(
        {
          subaccountId: [subaccountId],
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

    return {
      address,
      subaccountNumber,
      value: subaccountResponse.equity,
    };
  }

  @Get('/:address/:subaccountNumber/equity-list')
  public async getEquityList(
    @Path() address: string,
    @Path() subaccountNumber: number,
    @Query() createdOnOrAfter?: IsoString,
    @Query() createdBeforeOrAt?: IsoString,
  ): Promise<EquityListResponse> {
    const subaccountId: string = SubaccountTable.uuid(address, subaccountNumber);
    const subaccount: SubaccountFromDatabase | undefined = await SubaccountTable.findById(subaccountId);

    if (subaccount === undefined) {
      throw new NotFoundError(`No subaccount found with address ${address} and subaccountNumber ${subaccountNumber}`);
    }

    // Get PnL ticks in date range
    // Build query config conditionally to only include defined date filters
    const requiredFields: QueryableField[] = [];
    const queryConfig: {
      subaccountId: string[];
      createdOnOrAfter?: IsoString;
      createdBeforeOrAt?: IsoString;
    } = {
      subaccountId: [subaccountId],
    };

    if (createdOnOrAfter) {
      requiredFields.push(QueryableField.CREATED_ON_OR_AFTER);
      queryConfig.createdOnOrAfter = createdOnOrAfter;
    }
    if (createdBeforeOrAt) {
      requiredFields.push(QueryableField.CREATED_BEFORE_OR_AT);
      queryConfig.createdBeforeOrAt = createdBeforeOrAt;
    }

    // Query PnL ticks ordered by block height ascending (chronological order)
    const { results: pnlTicks } = await PnlTicksTable.findAll(
      queryConfig,
      requiredFields,
      {
        ...DEFAULT_POSTGRES_OPTIONS,
        orderBy: [[QueryableField.BLOCK_HEIGHT, Ordering.ASC]],
      },
    );

    // Map PnL ticks to equity list format
    // Use blockTime as the date since it represents when the equity was actually calculated
    const equityList = pnlTicks.map((tick: PnlTicksFromDatabase) => ({
      date: tick.blockTime,
      value: tick.equity,
    }));

    return {
      address,
      subaccountNumber,
      equityList,
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

router.get(
  '/:address/:subaccountNumber/live-pnl',
  rateLimiterMiddleware(getReqRateLimiter),
  ...CheckSubaccountSchema,
  ...CheckLimitAndCreatedBeforeOrAtAndOnOrAfterSchema,
  handleValidationErrors,
  complianceAndGeoCheck,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      address,
      subaccountNumber,
      createdOnOrAfter,
      createdBeforeOrAt,
    }: {
      address: string,
      subaccountNumber: number,
      createdOnOrAfter?: IsoString,
      createdBeforeOrAt?: IsoString,
    } = matchedData(req) as { address: string, subaccountNumber: number, createdOnOrAfter?: IsoString, createdBeforeOrAt?: IsoString };

    try {
      const controller: PortfolioController = new PortfolioController();
      const response: LivePnlResponse = await controller.getLivePnl(
        address,
        subaccountNumber,
        createdOnOrAfter,
        createdBeforeOrAt,
      );

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'PortfolioController GET /:address/:subaccountNumber/live-pnl',
        'Live PnL error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_live_pnl.timing`,
        Date.now() - start,
      );
    }
  },
);

router.get(
  '/:address/:subaccountNumber/realized-pnl',
  rateLimiterMiddleware(getReqRateLimiter),
  ...CheckSubaccountSchema,
  ...CheckLimitAndCreatedBeforeOrAtAndOnOrAfterSchema,
  handleValidationErrors,
  complianceAndGeoCheck,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      address,
      subaccountNumber,
      createdOnOrAfter,
      createdBeforeOrAt,
    }: {
      address: string,
      subaccountNumber: number,
      createdOnOrAfter?: IsoString,
      createdBeforeOrAt?: IsoString,
    } = matchedData(req) as { address: string, subaccountNumber: number, createdOnOrAfter?: IsoString, createdBeforeOrAt?: IsoString };

    try {
      const controller: PortfolioController = new PortfolioController();
      const response: RealizedPnlResponse = await controller.getRealizedPnl(
        address,
        subaccountNumber,
        createdOnOrAfter,
        createdBeforeOrAt,
      );

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'PortfolioController GET /:address/:subaccountNumber/realized-pnl',
        'Realized PnL error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_realized_pnl.timing`,
        Date.now() - start,
      );
    }
  },
);

router.get(
  '/:address/:subaccountNumber/profit-factor',
  rateLimiterMiddleware(getReqRateLimiter),
  ...CheckSubaccountSchema,
  ...CheckLimitAndCreatedBeforeOrAtAndOnOrAfterSchema,
  handleValidationErrors,
  complianceAndGeoCheck,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      address,
      subaccountNumber,
      createdOnOrAfter,
      createdBeforeOrAt,
    }: {
      address: string,
      subaccountNumber: number,
      createdOnOrAfter?: IsoString,
      createdBeforeOrAt?: IsoString,
    } = matchedData(req) as { address: string, subaccountNumber: number, createdOnOrAfter?: IsoString, createdBeforeOrAt?: IsoString };

    try {
      const controller: PortfolioController = new PortfolioController();
      const response: ProfitFactorResponse = await controller.getProfitFactor(
        address,
        subaccountNumber,
        createdOnOrAfter,
        createdBeforeOrAt,
      );

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'PortfolioController GET /:address/:subaccountNumber/profit-factor',
        'Profit factor error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_profit_factor.timing`,
        Date.now() - start,
      );
    }
  },
);

router.get(
  '/:address/:subaccountNumber/max-drawdown',
  rateLimiterMiddleware(getReqRateLimiter),
  ...CheckSubaccountSchema,
  ...CheckLimitAndCreatedBeforeOrAtAndOnOrAfterSchema,
  handleValidationErrors,
  complianceAndGeoCheck,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      address,
      subaccountNumber,
      createdOnOrAfter,
      createdBeforeOrAt,
    }: {
      address: string,
      subaccountNumber: number,
      createdOnOrAfter?: IsoString,
      createdBeforeOrAt?: IsoString,
    } = matchedData(req) as { address: string, subaccountNumber: number, createdOnOrAfter?: IsoString, createdBeforeOrAt?: IsoString };

    try {
      const controller: PortfolioController = new PortfolioController();
      const response: MaxDrawdownResponse = await controller.getMaxDrawdown(
        address,
        subaccountNumber,
        createdOnOrAfter,
        createdBeforeOrAt,
      );

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'PortfolioController GET /:address/:subaccountNumber/max-drawdown',
        'Max drawdown error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_max_drawdown.timing`,
        Date.now() - start,
      );
    }
  },
);

router.get(
  '/:address/:subaccountNumber/health',
  rateLimiterMiddleware(getReqRateLimiter),
  ...CheckSubaccountSchema,
  handleValidationErrors,
  complianceAndGeoCheck,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      address,
      subaccountNumber,
    }: {
      address: string,
      subaccountNumber: number,
    } = matchedData(req) as { address: string, subaccountNumber: number };

    try {
      const controller: PortfolioController = new PortfolioController();
      const response: HealthResponse = await controller.getHealth(
        address,
        subaccountNumber,
      );

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'PortfolioController GET /:address/:subaccountNumber/health',
        'Health error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_health.timing`,
        Date.now() - start,
      );
    }
  },
);

router.get(
  '/:address/:subaccountNumber/equity',
  rateLimiterMiddleware(getReqRateLimiter),
  ...CheckSubaccountSchema,
  handleValidationErrors,
  complianceAndGeoCheck,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      address,
      subaccountNumber,
    }: {
      address: string,
      subaccountNumber: number,
    } = matchedData(req) as { address: string, subaccountNumber: number };

    try {
      const controller: PortfolioController = new PortfolioController();
      const response: EquityResponse = await controller.getEquity(
        address,
        subaccountNumber,
      );

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'PortfolioController GET /:address/:subaccountNumber/equity',
        'Equity error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_equity.timing`,
        Date.now() - start,
      );
    }
  },
);

router.get(
  '/:address/:subaccountNumber/equity-list',
  rateLimiterMiddleware(getReqRateLimiter),
  ...CheckSubaccountSchema,
  ...CheckLimitAndCreatedBeforeOrAtAndOnOrAfterSchema,
  handleValidationErrors,
  complianceAndGeoCheck,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      address,
      subaccountNumber,
      createdOnOrAfter,
      createdBeforeOrAt,
    }: {
      address: string,
      subaccountNumber: number,
      createdOnOrAfter?: IsoString,
      createdBeforeOrAt?: IsoString,
    } = matchedData(req) as { address: string, subaccountNumber: number, createdOnOrAfter?: IsoString, createdBeforeOrAt?: IsoString };

    try {
      const controller: PortfolioController = new PortfolioController();
      const response: EquityListResponse = await controller.getEquityList(
        address,
        subaccountNumber,
        createdOnOrAfter,
        createdBeforeOrAt,
      );

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'PortfolioController GET /:address/:subaccountNumber/equity-list',
        'Equity list error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_equity_list.timing`,
        Date.now() - start,
      );
    }
  },
);

export default router;

