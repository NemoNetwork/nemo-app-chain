import { stats } from '@nemo-network-indexer/base/build';
import {
  SubaccountTable,
  IsoString,
  perpetualMarketRefresher,
  PerpetualMarketFromDatabase,
  FillTable,
  FillFromDatabase,
  QueryableField,
  FillColumns,
  Ordering,
  PerpetualPositionTable,
  PerpetualPositionFromDatabase,
  PositionSide,
  OrderSide,
} from '@nemo-network-indexer/postgres/build/src';
import Big from 'big.js';
import express from 'express';
import {
  checkSchema,
  matchedData,
  query,
} from 'express-validator';
import _ from 'lodash';
import {
  Controller, Get, Query, Route,
} from 'tsoa';

import { getReqRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import { complianceAndGeoCheck } from '../../../lib/compliance-and-geo-check';
import { NotFoundError } from '../../../lib/errors';
import {
  getChildSubaccountNums,
  getClobPairId, handleControllerError, isDefined,
} from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import {
  CheckLimitAndCreatedBeforeOrAtSchema,
  CheckSubaccountSchema,
  CheckParentSubaccountSchema,
  CheckPaginationSchema,
} from '../../../lib/validation/schemas';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import { fillToResponseObject } from '../../../request-helpers/request-transformer';
import {
  FillRequest,
  FillResponse,
  FillResponseObject,
  MarketAndTypeByClobPairId,
  MarketType,
  ParentSubaccountFillRequest,
} from '../../../types';

const router: express.Router = express.Router();
const controllerName: string = 'fills-controller';

@Route('fills')
class FillsController extends Controller {
  @Get('/')
  async getFills(
    @Query() address: string,
      @Query() subaccountNumber: number,
      @Query() market?: string,
      @Query() marketType?: MarketType,
      @Query() limit?: number,
      @Query() createdBeforeOrAtHeight?: number,
      @Query() createdBeforeOrAt?: IsoString,
      @Query() page?: number,
  ): Promise<FillResponse> {
    // TODO(DEC-656): Change to using a cache of markets in Redis similar to Librarian instead of
    // querying the DB.
    let clobPairId: string | undefined;
    if (isDefined(market) && isDefined(marketType)) {
      clobPairId = await getClobPairId(market!, marketType!);

      if (clobPairId === undefined) {
        throw new NotFoundError(`${market} not found in markets of type ${marketType}`);
      }
    }

    const subaccountId: string = SubaccountTable.uuid(address, subaccountNumber);
    const {
      results: fills,
      limit: pageSize,
      offset,
      total,
    } = await FillTable.findAll(
      {
        subaccountId: [subaccountId],
        clobPairId,
        limit,
        createdBeforeOrAtHeight: createdBeforeOrAtHeight
          ? createdBeforeOrAtHeight.toString()
          : undefined,
        createdBeforeOrAt,
        page,
      },
      [QueryableField.LIMIT],
      page !== undefined ? { orderBy: [[FillColumns.eventId, Ordering.ASC]] } : undefined,
    );

    const clobPairIdToPerpetualMarket: Record<
      string,
      PerpetualMarketFromDatabase> = perpetualMarketRefresher.getClobPairIdToPerpetualMarket();
    const clobPairIdToMarket: MarketAndTypeByClobPairId = _.mapValues(
      clobPairIdToPerpetualMarket,
      (perpetualMarket: PerpetualMarketFromDatabase) => {
        return {
          marketType: MarketType.PERPETUAL,
          market: perpetualMarket.ticker,
        };
      },
    );

    return {
      fills: fills.map((fill: FillFromDatabase): FillResponseObject => {
        return fillToResponseObject(fill, clobPairIdToMarket, subaccountNumber);
      }),
      pageSize,
      totalResults: total,
      offset,
    };
  }

  @Get('/parentSubaccount')
  async getFillsForParentSubaccount(
    @Query() address: string,
      @Query() parentSubaccountNumber: number,
      @Query() market?: string,
      @Query() marketType?: MarketType,
      @Query() limit?: number,
      @Query() createdBeforeOrAtHeight?: number,
      @Query() createdBeforeOrAt?: IsoString,
      @Query() page?: number,
  ): Promise<FillResponse> {
    // TODO(DEC-656): Change to using a cache of markets in Redis similar to Librarian instead of
    // querying the DB.
    let clobPairId: string | undefined;
    if (isDefined(market) && isDefined(marketType)) {
      clobPairId = await getClobPairId(market!, marketType!);

      if (clobPairId === undefined) {
        throw new NotFoundError(`${market} not found in markets of type ${marketType}`);
      }
    }

    // Get subaccountIds for all child subaccounts of the parent subaccount
    // Create a record of subaccountId to subaccount number
    const childIdtoSubaccountNumber: Record<string, number> = {};
    getChildSubaccountNums(parentSubaccountNumber).forEach(
      (subaccountNum: number) => {
        childIdtoSubaccountNumber[SubaccountTable.uuid(address, subaccountNum)] = subaccountNum;
      },
    );
    const subaccountIds: string[] = Object.keys(childIdtoSubaccountNumber);

    const {
      results: fills,
      limit: pageSize,
      offset,
      total,
    } = await FillTable.findAll(
      {
        subaccountId: subaccountIds,
        clobPairId,
        limit,
        createdBeforeOrAtHeight: createdBeforeOrAtHeight
          ? createdBeforeOrAtHeight.toString()
          : undefined,
        createdBeforeOrAt,
        page,
      },
      [QueryableField.LIMIT],
      page !== undefined ? { orderBy: [[FillColumns.eventId, Ordering.ASC]] } : undefined,
    );

    const clobPairIdToPerpetualMarket: Record<
        string,
        PerpetualMarketFromDatabase> = perpetualMarketRefresher.getClobPairIdToPerpetualMarket();
    const clobPairIdToMarket: MarketAndTypeByClobPairId = _.mapValues(
      clobPairIdToPerpetualMarket,
      (perpetualMarket: PerpetualMarketFromDatabase) => {
        return {
          marketType: MarketType.PERPETUAL,
          market: perpetualMarket.ticker,
        };
      },
    );

    // Query positions for all subaccounts to calculate closed PnL
    const positions: PerpetualPositionFromDatabase[] = await PerpetualPositionTable.findAll(
      {
        subaccountId: subaccountIds,
      },
      [],
    );

    // Create a map of (subaccountId, clobPairId) -> position for quick lookup
    // Use perpetualId from the position to match with clobPairId
    const positionMap: Record<string, PerpetualPositionFromDatabase> = {};
    positions.forEach((position: PerpetualPositionFromDatabase) => {
      const perpetualMarket = perpetualMarketRefresher.getPerpetualMarketFromId(position.perpetualId);
      if (perpetualMarket) {
        const key = `${position.subaccountId}-${perpetualMarket.clobPairId}`;
        positionMap[key] = position;
      }
    });

    // Helper function to calculate closed PnL for a fill
    const calculateClosedPnL = (fill: FillFromDatabase): string | undefined => {
      const perpetualMarket = clobPairIdToPerpetualMarket[fill.clobPairId];
      if (!perpetualMarket) {
        return undefined;
      }

      const key = `${fill.subaccountId}-${fill.clobPairId}`;
      const position = positionMap[key];

      if (!position) {
        return undefined;
      }

      // Check if fill closes the position (opposite sides)
      // A LONG position is closed by SELL fills, a SHORT position is closed by BUY fills
      const isClosingLong = position.side === PositionSide.LONG && fill.side === OrderSide.SELL;
      const isClosingShort = position.side === PositionSide.SHORT && fill.side === OrderSide.BUY;

      if (!isClosingLong && !isClosingShort) {
        return undefined;
      }

      // Calculate closed PnL
      // For LONG positions being closed by SELL: (fillPrice - entryPrice) * fillSize - fee
      // For SHORT positions being closed by BUY: (entryPrice - fillPrice) * fillSize - fee
      const fillPrice = Big(fill.price);
      const fillSize = Big(fill.size);
      const entryPrice = Big(position.entryPrice);
      const fee = Big(fill.fee || '0');

      let closedPnL: Big;
      if (isClosingLong) {
        // Long position: profit when selling above entry price
        closedPnL = fillPrice.minus(entryPrice).mul(fillSize).minus(fee);
      } else {
        // Short position: profit when buying below entry price
        closedPnL = entryPrice.minus(fillPrice).mul(fillSize).minus(fee);
      }

      return closedPnL.toFixed();
    };

    return {
      fills: fills.map((fill: FillFromDatabase): FillResponseObject => {
        const closedPnL = calculateClosedPnL(fill);
        return fillToResponseObject(fill, clobPairIdToMarket,
          childIdtoSubaccountNumber[fill.subaccountId], closedPnL);
      }),
      pageSize,
      totalResults: total,
      offset,
    };
  }
}

router.get(
  '/',
  rateLimiterMiddleware(getReqRateLimiter),
  ...CheckSubaccountSchema,
  ...CheckLimitAndCreatedBeforeOrAtSchema,
  ...CheckPaginationSchema,
  // Use conditional validations such that market is required if marketType is in the query
  // parameters and vice-versa.
  // Reference https://express-validator.github.io/docs/validation-chain-api.html#ifcondition
  query('market').if(query('marketType').exists()).notEmpty()
    .withMessage('market must be provided if marketType is provided'),
  query('marketType').if(query('market').exists()).notEmpty()
    .withMessage('marketType must be provided if market is provided'),
  // TODO(DEC-656): Validate market/marketType against cached markets.
  ...checkSchema({
    market: {
      in: ['query'],
      isString: true,
      optional: true,
    },
    marketType: {
      in: ['query'],
      isIn: {
        options: [Object.values(MarketType)],
      },
      optional: true,
      errorMessage: 'marketType must be a valid market type (PERPETUAL/SPOT)',
    },
  }),
  handleValidationErrors,
  complianceAndGeoCheck,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      address,
      subaccountNumber,
      market,
      marketType,
      limit,
      createdBeforeOrAtHeight,
      createdBeforeOrAt,
      page,
    }: FillRequest = matchedData(req) as FillRequest;

    // The schema checks allow subaccountNumber to be a string, but we know it's a number here.
    const subaccountNum : number = +subaccountNumber;

    // TODO(DEC-656): Change to using a cache of markets in Redis similar to Librarian instead of
    // querying the DB.
    try {
      const controller: FillsController = new FillsController();
      const response: FillResponse = await controller.getFills(
        address,
        subaccountNum,
        market,
        marketType,
        limit,
        createdBeforeOrAtHeight,
        createdBeforeOrAt,
        page,
      );

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'FillsController GET /',
        'Fills error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_fills.timing`,
        Date.now() - start,
      );
    }
  },
);

router.get(
  '/parentSubaccountNumber',
  rateLimiterMiddleware(getReqRateLimiter),
  ...CheckParentSubaccountSchema,
  ...CheckLimitAndCreatedBeforeOrAtSchema,
  ...CheckPaginationSchema,
  // Use conditional validations such that market is required if marketType is in the query
  // parameters and vice-versa.
  // Reference https://express-validator.github.io/docs/validation-chain-api.html#ifcondition
  query('market').if(query('marketType').exists()).notEmpty()
    .withMessage('market must be provided if marketType is provided'),
  query('marketType').if(query('market').exists()).notEmpty()
    .withMessage('marketType must be provided if market is provided'),
  // TODO(DEC-656): Validate market/marketType against cached markets.
  ...checkSchema({
    market: {
      in: ['query'],
      isString: true,
      optional: true,
    },
    marketType: {
      in: ['query'],
      isIn: {
        options: [Object.values(MarketType)],
      },
      optional: true,
      errorMessage: 'marketType must be a valid market type (PERPETUAL/SPOT)',
    },
  }),
  handleValidationErrors,
  complianceAndGeoCheck,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      address,
      parentSubaccountNumber,
      market,
      marketType,
      limit,
      createdBeforeOrAtHeight,
      createdBeforeOrAt,
      page,
    }: ParentSubaccountFillRequest = matchedData(req) as ParentSubaccountFillRequest;

    // The schema checks allow subaccountNumber to be a string, but we know it's a number here.
    const parentSubaccountNum : number = +parentSubaccountNumber;

    // TODO(DEC-656): Change to using a cache of markets in Redis similar to Librarian instead of
    // querying the DB.
    try {
      const controller: FillsController = new FillsController();
      const response: FillResponse = await controller.getFillsForParentSubaccount(
        address,
        parentSubaccountNum,
        market,
        marketType,
        limit,
        createdBeforeOrAtHeight,
        createdBeforeOrAt,
        page,
      );

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'FillsController GET /parentSubaccountNumber',
        'Fills error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_fills.timing`,
        Date.now() - start,
      );
    }
  },
);

export default router;
