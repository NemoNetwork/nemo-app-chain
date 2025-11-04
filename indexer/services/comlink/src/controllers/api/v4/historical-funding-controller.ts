import { stats } from '@nemo-network-indexer/base/build';
import {
  DEFAULT_POSTGRES_OPTIONS,
  FillTable,
  FundingIndexUpdatesColumns,
  FundingIndexUpdatesFromDatabase,
  FundingIndexUpdatesTable,
  IsoString,
  Ordering,
  PerpetualMarketFromDatabase,
  PerpetualMarketTable,
  SubaccountTable,
  SubaccountFromDatabase,
  PerpetualPositionTable,
  PerpetualPositionFromDatabase,
  PerpetualPositionStatus,
  PositionSide,
  perpetualMarketRefresher,
} from '@nemo-network-indexer/postgres/build/src';
import Big from 'big.js';
import express from 'express';
import { matchedData } from 'express-validator';
import {
  Controller, Get, Path, Query, Route,
} from 'tsoa';

import { getReqRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import { NotFoundError } from '../../../lib/errors';
import { handleControllerError } from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import { CheckEffectiveBeforeOrAtSchema, CheckLimitSchema, CheckTickerParamSchema, CheckSubaccountHistoricalFundingSchema } from '../../../lib/validation/schemas';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import { historicalFundingToResponseObject } from '../../../request-helpers/request-transformer';
import { quantumsToHuman } from '@nemo-network-indexer/postgres/build/src/lib/protocol-translations';
import {
  HistoricalFundingRequest,
  HistoricalFundingResponse,
  HistoricalFundingResponseObject,
  MarketType,
  SubaccountHistoricalFundingRequest,
  SubaccountHistoricalFundingResponse,
  SubaccountHistoricalFundingResponseObject,
} from '../../../types';

const router: express.Router = express.Router();
const controllerName: string = 'historical-funding-controller';

@Route('historicalFunding')
class HistoricalFundingController extends Controller {
  @Get('/:ticker')
  async getHistoricalFunding(
    @Path() ticker: string,
      @Query() limit?: number,
      @Query() effectiveBeforeOrAtHeight?: number,
      @Query() effectiveBeforeOrAt?: IsoString,
  ): Promise<HistoricalFundingResponse> {
    const perpetualMarket: (
      PerpetualMarketFromDatabase | undefined
    ) = await PerpetualMarketTable.findByTicker(ticker);

    if (perpetualMarket === undefined) {
      throw new NotFoundError(`${ticker} not found in markets of type ${MarketType.PERPETUAL}`);
    }

    const fundingIndices: FundingIndexUpdatesFromDatabase[] = await
    FundingIndexUpdatesTable.findAll(
      {
        perpetualId: [perpetualMarket.id],
        effectiveBeforeOrAt,
        effectiveBeforeOrAtHeight: effectiveBeforeOrAtHeight
          ? effectiveBeforeOrAtHeight.toString()
          : undefined,
        limit,
      }, [],
      {
        ...DEFAULT_POSTGRES_OPTIONS,
        orderBy: [[FundingIndexUpdatesColumns.effectiveAtHeight, Ordering.DESC]],
      },
    );

    return {
      historicalFunding: fundingIndices.map(
        (fundingIndex: FundingIndexUpdatesFromDatabase): HistoricalFundingResponseObject => {
          return historicalFundingToResponseObject(fundingIndex, ticker);
        },
      ),
    };
  }

  @Get('/address/:address/:subaccountNumber')
  async getSubaccountHistoricalFunding(
    @Path() address: string,
    @Path() subaccountNumber: number,
    @Query() limit?: number,
    @Query() effectiveBeforeOrAtHeight?: number,
    @Query() effectiveBeforeOrAt?: IsoString,
  ): Promise<SubaccountHistoricalFundingResponse> {
    // Get the subaccount
    const subaccountId: string = SubaccountTable.uuid(address, subaccountNumber);
    const subaccount: SubaccountFromDatabase | undefined = await SubaccountTable.findById(subaccountId);

    if (subaccount === undefined) {
      throw new NotFoundError(`No subaccount found with address ${address} and subaccountNumber ${subaccountNumber}`);
    }

    // Get all perpetual positions for this subaccount
    const perpetualPositions: PerpetualPositionFromDatabase[] = await PerpetualPositionTable.findAll(
      {
        subaccountId: [subaccountId],
        status: [PerpetualPositionStatus.OPEN],
      },
      [],
      {
        ...DEFAULT_POSTGRES_OPTIONS,
      },
    );

    if (perpetualPositions.length === 0) {
      return {
        historicalFunding: [],
      };
    }

    // Get perpetual market data for all positions
    const perpetualIds = perpetualPositions.map(pos => pos.perpetualId);
    const perpetualMarkets: PerpetualMarketFromDatabase[] = await PerpetualMarketTable.findAll(
      {
        id: perpetualIds,
      },
      [],
      {
        ...DEFAULT_POSTGRES_OPTIONS,
      },
    );

    const perpetualMarketsMap: { [id: string]: PerpetualMarketFromDatabase } = perpetualMarkets.reduce((acc, market) => {
      acc[market.id.toString()] = market;
      return acc;
    }, {} as { [id: string]: PerpetualMarketFromDatabase });

    // Get funding index updates for all perpetual markets, ordered by height DESC
    // We need to reverse to ASC to calculate payment from previous index
    const fundingIndices: FundingIndexUpdatesFromDatabase[] = await FundingIndexUpdatesTable.findAll(
      {
        perpetualId: perpetualIds,
        effectiveBeforeOrAt,
        effectiveBeforeOrAtHeight: effectiveBeforeOrAtHeight
          ? effectiveBeforeOrAtHeight.toString()
          : undefined,
        limit,
      }, [],
      {
        ...DEFAULT_POSTGRES_OPTIONS,
        orderBy: [[FundingIndexUpdatesColumns.effectiveAtHeight, Ordering.DESC]],
      },
    );

    // Reverse to ascending order for calculating payment from funding index changes
    const fundingIndicesAsc = [...fundingIndices].reverse();

    // Create a map to track previous funding index for each perpetual
    const previousFundingIndexMap: { [perpetualId: string]: string | undefined } = {};
    // Create a map to track historical position sizes at each funding index update
    const historicalPositionSizeMap: { [perpetualId: string]: { [height: string]: string } } = {};

    // Pre-calculate historical position sizes for each funding index update
    await perpetualMarketRefresher.updatePerpetualMarkets();
    const clobPairIdToPerpetualId: { [clobPairId: string]: string } = {};
    for (const market of perpetualMarkets) {
      // Use the market's clobPairId directly from the database
      clobPairIdToPerpetualId[market.clobPairId] = market.id.toString();
    }

    // For each unique effectiveAtHeight, get historical position sizes
    const uniqueHeights = [...new Set(fundingIndicesAsc.map(fi => fi.effectiveAtHeight))];
    for (const height of uniqueHeights) {
      const openSizes = await FillTable.getOpenSizeWithFundingIndex(subaccountId, height);
      for (const openSize of openSizes) {
        const perpetualId = clobPairIdToPerpetualId[openSize.clobPairId];
        if (perpetualId) {
          const market = perpetualMarkets.find(m => m.id.toString() === perpetualId);
          if (market) {
            // Convert openSize from base quantums to human-readable format
            // openSize from getOpenSizeWithFundingIndex is in base quantums (signed)
            // Use quantumsToHuman helper to convert to human-readable format
            // Note: quantumsToHuman multiplies by 10^atomicResolution, preserving the sign
            const openSizeHuman = quantumsToHuman(openSize.openSize, market.atomicResolution).toFixed();
            
            if (!historicalPositionSizeMap[perpetualId]) {
              historicalPositionSizeMap[perpetualId] = {};
            }
            // Store the converted position size (preserves sign: positive = LONG, negative = SHORT)
            historicalPositionSizeMap[perpetualId][height] = openSizeHuman;
          }
        }
      }
    }

    // Create response objects
    const historicalFunding: SubaccountHistoricalFundingResponseObject[] = fundingIndicesAsc.map(
      (fundingIndex: FundingIndexUpdatesFromDatabase): SubaccountHistoricalFundingResponseObject => {
        const perpetualMarket = perpetualMarketsMap[fundingIndex.perpetualId];
        const position = perpetualPositions.find(pos => pos.perpetualId === fundingIndex.perpetualId);
        
        // Get historical position size at this funding index update time
        const historicalPositionSize = historicalPositionSizeMap[fundingIndex.perpetualId]?.[fundingIndex.effectiveAtHeight];
        // Use historical position size if available, otherwise fall back to current position size
        // If no position exists at all, use '0' (no position)
        const positionSizeStr = historicalPositionSize !== undefined 
          ? historicalPositionSize 
          : (position ? position.size : '0');
        const positionSize = Big(positionSizeStr);
        
        // Determine position type from position size sign (positive = LONG, negative = SHORT)
        // Note: positionSize of exactly 0 means no position, we'll show as LONG (could also be SHORT, but 0 is edge case)
        const positionType = positionSize.gt(0) ? 'LONG' : positionSize.lt(0) ? 'SHORT' : 'LONG';
        
        let payment = '0';
        if (position && !positionSize.eq(0)) {
          const currentFundingIndex = Big(fundingIndex.fundingIndex);
          const previousFundingIndex = previousFundingIndexMap[fundingIndex.perpetualId];
          
          if (previousFundingIndex !== undefined) {
            // Calculate payment based on funding index change
            // Using the same convention as getUnsettledFunding:
            // Payment = positionSize * (previousFundingIndex - currentFundingIndex)
            // This ensures correct signs: when funding index increases, shorts get paid (positive)
            const fundingIndexChange = Big(previousFundingIndex).minus(currentFundingIndex);
            payment = positionSize.times(fundingIndexChange).toFixed();
          } else {
            // For the first funding index update, we can't calculate payment
            // Use settledFunding as fallback, but this is cumulative
            payment = position.settledFunding;
          }
          
          // Update previous funding index for next iteration
          previousFundingIndexMap[fundingIndex.perpetualId] = fundingIndex.fundingIndex;
        }
        
        return {
          market: perpetualMarket?.ticker || 'UNKNOWN',
          positionType,
          date: fundingIndex.effectiveAt,
          positionSize: positionSizeStr,
          payment,
          fundingRate: fundingIndex.rate,
        };
      },
    );

    // Reverse back to DESC order for response
    historicalFunding.reverse();

    return {
      historicalFunding,
    };
  }
}

router.get(
  '/:ticker',
  rateLimiterMiddleware(getReqRateLimiter),
  ...CheckLimitSchema,
  ...CheckTickerParamSchema,
  ...CheckEffectiveBeforeOrAtSchema,
  handleValidationErrors,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      ticker,
      limit,
      effectiveBeforeOrAtHeight,
      effectiveBeforeOrAt,
    }: {
      ticker: string,
      limit: number,
      effectiveBeforeOrAtHeight?: number,
      effectiveBeforeOrAt?: IsoString,
    } = matchedData(req) as HistoricalFundingRequest;

    try {
      const controller: HistoricalFundingController = new HistoricalFundingController();
      const response: HistoricalFundingResponse = await controller.getHistoricalFunding(
        ticker,
        limit,
        effectiveBeforeOrAtHeight,
        effectiveBeforeOrAt,
      );

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'HistoricalFundingController GET /',
        'HistoricalFunding error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_historical_funding.timing`,
        Date.now() - start,
      );
    }
  },
);

router.get(
  '/address/:address/:subaccountNumber',
  rateLimiterMiddleware(getReqRateLimiter),
  ...CheckSubaccountHistoricalFundingSchema,
  handleValidationErrors,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      address,
      subaccountNumber,
      limit,
      effectiveBeforeOrAtHeight,
      effectiveBeforeOrAt,
    }: {
      address: string,
      subaccountNumber: number,
      limit: number,
      effectiveBeforeOrAtHeight?: number,
      effectiveBeforeOrAt?: IsoString,
    } = matchedData(req) as SubaccountHistoricalFundingRequest;

    try {
      const controller: HistoricalFundingController = new HistoricalFundingController();
      const response: SubaccountHistoricalFundingResponse = await controller.getSubaccountHistoricalFunding(
        address,
        subaccountNumber,
        limit,
        effectiveBeforeOrAtHeight,
        effectiveBeforeOrAt,
      );

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'HistoricalFundingController GET /subaccount/:address/:subaccountNumber',
        'Subaccount historical funding error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_subaccount_historical_funding.timing`,
        Date.now() - start,
      );
    }
  },
);

export default router;
