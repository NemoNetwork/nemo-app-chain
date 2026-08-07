import { stats } from '@nemo-network-indexer/base/build';
import {
  PnlTicksFromDatabase,
  perpetualMarketRefresher,
  PerpetualMarketFromDatabase,
  USDC_ASSET_ID,
  FundingIndexMap,
  AssetPositionFromDatabase,
  PerpetualPositionFromDatabase,
  SubaccountFromDatabase,
  AssetColumns,
  BlockTable,
  MarketTable,
  AssetPositionTable,
  PerpetualPositionStatus,
  PerpetualPositionTable,
  AssetTable,
  SubaccountTable,
  AssetFromDatabase,
  MarketFromDatabase,
  BlockFromDatabase,
  FillTable,
  FundingIndexUpdatesTable,
  IsoString,
  PnlTickInterval,
  VaultTable,
  VaultFromDatabase,
  MEGAVAULT_SUBACCOUNT_ID,
  TransferFromDatabase,
  TransferTable,
  TransferColumns,
  Ordering,
  VaultPnlTicksView,
} from '@nemo-network-indexer/postgres/build/src';
import Big from 'big.js';
import bounds from 'binary-searching';
import express from 'express';
import { checkSchema, matchedData } from 'express-validator';
import _, { Dictionary } from 'lodash';
import { DateTime } from 'luxon';
import {
  Controller, Get, Query, Route,
} from 'tsoa';

import { getReqRateLimiter } from '../../../caches/rate-limiters';
import { getVaultStartPnl } from '../../../caches/vault-start-pnl';
import config from '../../../config';
import {
  aggregateHourlyPnlTicks,
  getSubaccountResponse,
  getVaultMapping,
  getVaultPnlStartDate,
  handleControllerError,
} from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import { CheckLimitAndCreatedBeforeOrAtSchema } from '../../../lib/validation/schemas';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import { pnlTicksToResponseObject } from '../../../request-helpers/request-transformer';
import {
  MegavaultHistoricalPnlResponse,
  VaultsHistoricalPnlResponse,
  VaultHistoricalPnl,
  VaultPosition,
  AssetById,
  MegavaultPositionResponse,
  SubaccountResponseObject,
  MegavaultHistoricalPnlRequest,
  VaultsHistoricalPnlRequest,
  AggregatedPnlTick,
  PnlTicksResponseObject,
  VaultMapping,
  VaultsResponse,
  VaultResponseObject,
  MegavaultSummaryResponse,
  MegavaultTransfersResponse,
  MegavaultTransfersRequest,
  MegavaultTransferResponseObject,
  MegavaultTransferType,
  MegavaultTransferStatus,
  MegavaultTransferStatusResponse,
  MegavaultTransferStatusRequest,
} from '../../../types';

const router: express.Router = express.Router();
const controllerName: string = 'vault-controller';

@Route('vault/v1')
class VaultController extends Controller {
  @Get('/megavault/historicalPnl')
  async getMegavaultHistoricalPnl(
    @Query() resolution?: PnlTickInterval,
  ): Promise<MegavaultHistoricalPnlResponse> {
    const start: number = Date.now();
    const vaultSubaccounts: VaultMapping = await getVaultMapping();
    stats.timing(
      `${config.SERVICE_NAME}.${controllerName}.fetch_vaults.timing`,
      Date.now() - start,
    );

    const startTicksPositions: number = Date.now();
    const vaultSubaccountIdsWithMainSubaccount: string[] = _
      .keys(vaultSubaccounts)
      .concat([MEGAVAULT_SUBACCOUNT_ID]);
    const [
      vaultPnlTicks,
      vaultPositions,
      latestBlock,
      mainSubaccountEquity,
      latestPnlTick,
      firstMainVaultTransferTimestamp,
    ] : [
      PnlTicksFromDatabase[],
      Map<string, VaultPosition>,
      BlockFromDatabase,
      string,
      PnlTicksFromDatabase | undefined,
      DateTime | undefined,
    ] = await Promise.all([
      getVaultSubaccountPnlTicks(vaultSubaccountIdsWithMainSubaccount, getResolution(resolution)),
      getVaultPositions(vaultSubaccounts),
      BlockTable.getLatest(),
      getMainSubaccountEquity(),
      getLatestPnlTick(_.values(vaultSubaccounts)),
      getFirstMainVaultTransferDateTime(),
    ]);
    stats.timing(
      `${config.SERVICE_NAME}.${controllerName}.fetch_ticks_positions_equity.timing`,
      Date.now() - startTicksPositions,
    );
    // aggregate pnlTicks for all vault subaccounts grouped by blockHeight
    const aggregatedPnlTicks: PnlTicksFromDatabase[] = aggregateVaultPnlTicks(
      vaultPnlTicks,
      _.values(vaultSubaccounts),
      firstMainVaultTransferTimestamp,
    );

    const currentEquity: string = Array.from(vaultPositions.values())
      .map((position: VaultPosition): string => {
        return position.equity;
      }).reduce((acc: string, curr: string): string => {
        return (Big(acc).add(Big(curr))).toFixed();
      }, mainSubaccountEquity);
    const pnlTicksWithCurrentTick: PnlTicksFromDatabase[] = getPnlTicksWithCurrentTick(
      currentEquity,
      filterOutIntervalTicks(aggregatedPnlTicks, getResolution(resolution)),
      latestBlock,
      latestPnlTick,
    );

    return {
      megavaultPnl: _.sortBy(pnlTicksWithCurrentTick, 'blockTime').map(
        (pnlTick: PnlTicksFromDatabase) => {
          return pnlTicksToResponseObject(pnlTick);
        }),
    };
  }

  @Get('/vaults/historicalPnl')
  async getVaultsHistoricalPnl(
    @Query() resolution?: PnlTickInterval,
  ): Promise<VaultsHistoricalPnlResponse> {
    const vaultSubaccounts: VaultMapping = await getVaultMapping();
    const [
      vaultPnlTicks,
      vaultPositions,
      latestBlock,
      latestTicks,
    ] : [
      PnlTicksFromDatabase[],
      Map<string, VaultPosition>,
      BlockFromDatabase,
      PnlTicksFromDatabase[],
    ] = await Promise.all([
      getVaultSubaccountPnlTicks(_.keys(vaultSubaccounts), getResolution(resolution)),
      getVaultPositions(vaultSubaccounts),
      BlockTable.getLatest(),
      getLatestPnlTicks(),
    ]);
    const latestTicksBySubaccountId: Dictionary<PnlTicksFromDatabase> = _.keyBy(
      latestTicks,
      'subaccountId',
    );

    const groupedVaultPnlTicks: VaultHistoricalPnl[] = _(vaultPnlTicks)
      .filter((pnlTickFromDatabsae: PnlTicksFromDatabase): boolean => {
        return vaultSubaccounts[pnlTickFromDatabsae.subaccountId] !== undefined;
      })
      .groupBy('subaccountId')
      .mapValues((pnlTicks: PnlTicksFromDatabase[], subaccountId: string): VaultHistoricalPnl => {
        const market: PerpetualMarketFromDatabase | undefined = perpetualMarketRefresher
          .getPerpetualMarketFromClobPairId(
            vaultSubaccounts[subaccountId].clobPairId,
          );

        if (market === undefined) {
          throw new Error(
            `Vault clob pair id ${vaultSubaccounts[subaccountId]} does not correspond to ` +
            'a perpetual market.');
        }

        const vaultPosition: VaultPosition | undefined = vaultPositions.get(subaccountId);
        const currentEquity: string = vaultPosition === undefined ? '0' : vaultPosition.equity;
        const pnlTicksWithCurrentTick: PnlTicksFromDatabase[] = getPnlTicksWithCurrentTick(
          currentEquity,
          pnlTicks,
          latestBlock,
          latestTicksBySubaccountId[subaccountId],
        );

        // Only retain fields we need and excludes fields like `id`, `subaccountId`.
        const pnlTicksResponseObjects
        : PnlTicksResponseObject[] = pnlTicksWithCurrentTick.map(pnlTicksToResponseObject);
        return {
          ticker: market.ticker,
          historicalPnl: pnlTicksResponseObjects,
        };
      })
      .values()
      .value();

    return {
      vaultsPnl: _.sortBy(groupedVaultPnlTicks, 'ticker'),
    };
  }

  @Get('/megavault/positions')
  async getMegavaultPositions(): Promise<MegavaultPositionResponse> {
    const vaultSubaccounts: VaultMapping = await getVaultMapping();

    const vaultPositions: Map<string, VaultPosition> = await getVaultPositions(vaultSubaccounts);

    return {
      positions: _.sortBy(Array.from(vaultPositions.values()), 'ticker'),
    };
  }

  @Get('/vaults')
  async getVaults(): Promise<VaultsResponse> {
    const vaults: VaultFromDatabase[] = await VaultTable.findAll({}, [], {});

    const vaultObjects: VaultResponseObject[] = vaults.reduce(
      (validVaults: VaultResponseObject[], vault: VaultFromDatabase): VaultResponseObject[] => {
        const perpetualMarket: PerpetualMarketFromDatabase | undefined = perpetualMarketRefresher
          .getPerpetualMarketFromClobPairId(vault.clobPairId);
        if (perpetualMarket === undefined) {
          return validVaults;
        }
        return validVaults.concat({
          address: vault.address,
          ticker: perpetualMarket.ticker,
          status: vault.status,
          createdAt: vault.createdAt,
          updatedAt: vault.updatedAt,
        });
      },
      [],
    );

    return {
      vaults: _.sortBy(vaultObjects, 'ticker'),
    };
  }

  @Get('/megavault/summary')
  async getMegavaultSummary(): Promise<MegavaultSummaryResponse> {
    const vaultSubaccounts: VaultMapping = await getVaultMapping();

    const [
      megavaultPnlResponse,
      vaultPositions,
      mainSubaccountEquity,
      volume24H,
      firstMainVaultTransferTimestamp,
    ] : [
      MegavaultHistoricalPnlResponse,
      Map<string, VaultPosition>,
      string,
      Big,
      DateTime | undefined,
    ] = await Promise.all([
      this.getMegavaultHistoricalPnl(PnlTickInterval.day),
      getVaultPositions(vaultSubaccounts),
      getMainSubaccountEquity(),
      FillTable.getTotalVolumeForSubaccounts(
        _.keys(vaultSubaccounts),
        DateTime.utc().minus({ hours: 24 }).toISO(),
      ),
      getFirstMainVaultTransferDateTime(),
    ]);

    const equity: string = Array.from(vaultPositions.values())
      .map((position: VaultPosition): string => {
        return position.equity;
      }).reduce((acc: string, curr: string): string => {
        return (Big(acc).add(Big(curr))).toFixed();
      }, mainSubaccountEquity);

    // Response ticks are already sorted by block time, ending with a synthetic current tick.
    const pnlTicks: PnlTicksResponseObject[] = megavaultPnlResponse.megavaultPnl;
    const latestTick: PnlTicksResponseObject | undefined = _.last(pnlTicks);

    return {
      equity,
      numVaults: _.keys(vaultSubaccounts).length,
      allTimePnl: latestTick === undefined ? '0' : latestTick.totalPnl,
      apr: computeAnnualizedReturn(pnlTicks),
      maxDrawdown: computeMaxPnlDrawdown(pnlTicks),
      volume24H: volume24H.toFixed(),
      createdAt: firstMainVaultTransferTimestamp === undefined
        ? null
        : firstMainVaultTransferTimestamp.toISO(),
    };
  }

  @Get('/megavault/transfers')
  async getMegavaultTransfers(
    @Query() address: string,
      @Query() limit?: number,
      @Query() createdBeforeOrAt?: IsoString,
      @Query() createdBeforeOrAtHeight?: number,
  ): Promise<MegavaultTransfersResponse> {
    const responseLimit: number = limit ?? config.API_LIMIT_V4;
    const userSubaccounts: SubaccountFromDatabase[] = await SubaccountTable.findAll(
      { address },
      [],
    );
    if (userSubaccounts.length === 0) {
      return { transfers: [] };
    }
    const userSubaccountIds: string[] = userSubaccounts.map(
      (subaccount: SubaccountFromDatabase): string => { return subaccount.id; },
    );

    const [
      deposits,
      withdrawals,
      assets,
    ] : [
      TransferFromDatabase[],
      TransferFromDatabase[],
      AssetFromDatabase[],
    ] = await Promise.all([
      TransferTable.findAll(
        {
          senderSubaccountId: userSubaccountIds,
          recipientSubaccountId: [MEGAVAULT_SUBACCOUNT_ID],
          createdBeforeOrAt,
          createdBeforeOrAtHeight: createdBeforeOrAtHeight !== undefined
            ? createdBeforeOrAtHeight.toString()
            : undefined,
          limit: responseLimit,
        },
        [],
        { orderBy: [[TransferColumns.createdAtHeight, Ordering.DESC]] },
      ),
      TransferTable.findAll(
        {
          senderSubaccountId: [MEGAVAULT_SUBACCOUNT_ID],
          recipientSubaccountId: userSubaccountIds,
          createdBeforeOrAt,
          createdBeforeOrAtHeight: createdBeforeOrAtHeight !== undefined
            ? createdBeforeOrAtHeight.toString()
            : undefined,
          limit: responseLimit,
        },
        [],
        { orderBy: [[TransferColumns.createdAtHeight, Ordering.DESC]] },
      ),
      AssetTable.findAll({}, []),
    ]);

    const assetMap: AssetById = _.keyBy(assets, AssetColumns.id);
    const subaccountsById: { [id: string]: SubaccountFromDatabase } = _.keyBy(
      userSubaccounts,
      'id',
    );
    const mergedTransfers: TransferFromDatabase[] = deposits.concat(withdrawals)
      .sort((a: TransferFromDatabase, b: TransferFromDatabase): number => {
        return parseInt(b.createdAtHeight, 10) - parseInt(a.createdAtHeight, 10);
      })
      .slice(0, responseLimit);

    return {
      transfers: mergedTransfers.map(
        (transfer: TransferFromDatabase): MegavaultTransferResponseObject => {
          return megavaultTransferToResponseObject(transfer, assetMap, subaccountsById);
        },
      ),
    };
  }

  @Get('/megavault/transfers/status')
  async getMegavaultTransferStatus(
    @Query() transactionHash: string,
  ): Promise<MegavaultTransferStatusResponse> {
    const [
      transfers,
      assets,
    ] : [
      TransferFromDatabase[],
      AssetFromDatabase[],
    ] = await Promise.all([
      TransferTable.findAll({ transactionHash: [transactionHash] }, []),
      AssetTable.findAll({}, []),
    ]);

    const megavaultTransfers: TransferFromDatabase[] = transfers.filter(
      (transfer: TransferFromDatabase): boolean => {
        return transfer.senderSubaccountId === MEGAVAULT_SUBACCOUNT_ID ||
          transfer.recipientSubaccountId === MEGAVAULT_SUBACCOUNT_ID;
      },
    );

    const counterpartySubaccountIds: string[] = _.uniq(
      megavaultTransfers
        .map((transfer: TransferFromDatabase): string | undefined => {
          return transfer.senderSubaccountId === MEGAVAULT_SUBACCOUNT_ID
            ? transfer.recipientSubaccountId
            : transfer.senderSubaccountId;
        })
        .filter((subaccountId: string | undefined): boolean => {
          return subaccountId !== undefined;
        }) as string[],
    );
    const counterpartySubaccounts: SubaccountFromDatabase[] = counterpartySubaccountIds.length > 0
      ? await SubaccountTable.findAll({ id: counterpartySubaccountIds }, [])
      : [];

    const assetMap: AssetById = _.keyBy(assets, AssetColumns.id);
    const subaccountsById: { [id: string]: SubaccountFromDatabase } = _.keyBy(
      counterpartySubaccounts,
      'id',
    );

    return {
      transactionHash,
      // A transaction that fails on-chain is never indexed, so the indexer can only ever
      // distinguish "seen" (COMPLETED) from "not yet seen" (PENDING). Clients are expected
      // to stop polling after their own timeout.
      status: megavaultTransfers.length > 0
        ? MegavaultTransferStatus.COMPLETED
        : MegavaultTransferStatus.PENDING,
      transfers: megavaultTransfers.map(
        (transfer: TransferFromDatabase): MegavaultTransferResponseObject => {
          return megavaultTransferToResponseObject(transfer, assetMap, subaccountsById);
        },
      ),
    };
  }
}

router.get(
  '/v1/megavault/historicalPnl',
  ...checkSchema({
    resolution: {
      in: 'query',
      isIn: {
        options: [Object.values(PnlTickInterval)],
        errorMessage: `type must be one of ${Object.values(PnlTickInterval)}`,
      },
      optional: true,
    },
  }),
  handleValidationErrors,
  rateLimiterMiddleware(getReqRateLimiter),
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      resolution,
    }: MegavaultHistoricalPnlRequest = matchedData(req) as MegavaultHistoricalPnlRequest;

    try {
      const controllers: VaultController = new VaultController();
      const response: MegavaultHistoricalPnlResponse = await controllers
        .getMegavaultHistoricalPnl(resolution);
      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'VaultController GET /megavault/historicalPnl',
        'Megavault Historical Pnl error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_megavault_historical_pnl.timing`,
        Date.now() - start,
      );
    }
  },
);

router.get(
  '/v1/vaults/historicalPnl',
  ...checkSchema({
    resolution: {
      in: 'query',
      isIn: {
        options: [Object.values(PnlTickInterval)],
        errorMessage: `type must be one of ${Object.values(PnlTickInterval)}`,
      },
      optional: true,
    },
  }),
  handleValidationErrors,
  rateLimiterMiddleware(getReqRateLimiter),
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      resolution,
    }: VaultsHistoricalPnlRequest = matchedData(req) as VaultsHistoricalPnlRequest;

    try {
      const controllers: VaultController = new VaultController();
      const response: VaultsHistoricalPnlResponse = await controllers
        .getVaultsHistoricalPnl(resolution);
      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'VaultHistoricalPnlController GET /vaults/historicalPnl',
        'Vaults Historical Pnl error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_vaults_historical_pnl.timing`,
        Date.now() - start,
      );
    }
  },
);

router.get(
  '/v1/megavault/positions',
  rateLimiterMiddleware(getReqRateLimiter),
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    try {
      const controllers: VaultController = new VaultController();
      const response: MegavaultPositionResponse = await controllers.getMegavaultPositions();
      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'VaultController GET /megavault/positions',
        'Megavault Positions error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_megavault_positions.timing`,
        Date.now() - start,
      );
    }
  },
);

router.get(
  '/v1/vaults',
  rateLimiterMiddleware(getReqRateLimiter),
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    try {
      const controllers: VaultController = new VaultController();
      const response: VaultsResponse = await controllers.getVaults();
      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'VaultController GET /vaults',
        'Vaults error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_vaults.timing`,
        Date.now() - start,
      );
    }
  },
);

router.get(
  '/v1/megavault/summary',
  rateLimiterMiddleware(getReqRateLimiter),
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    try {
      const controllers: VaultController = new VaultController();
      const response: MegavaultSummaryResponse = await controllers.getMegavaultSummary();
      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'VaultController GET /megavault/summary',
        'Megavault Summary error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_megavault_summary.timing`,
        Date.now() - start,
      );
    }
  },
);

router.get(
  '/v1/megavault/transfers',
  ...checkSchema({
    address: {
      in: 'query',
      isString: true,
      notEmpty: true,
      errorMessage: 'address must be a non-empty string',
    },
  }),
  ...CheckLimitAndCreatedBeforeOrAtSchema,
  handleValidationErrors,
  rateLimiterMiddleware(getReqRateLimiter),
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      address,
      limit,
      createdBeforeOrAt,
      createdBeforeOrAtHeight,
    }: MegavaultTransfersRequest = matchedData(req) as MegavaultTransfersRequest;

    try {
      const controllers: VaultController = new VaultController();
      const response: MegavaultTransfersResponse = await controllers.getMegavaultTransfers(
        address,
        limit,
        createdBeforeOrAt,
        createdBeforeOrAtHeight,
      );
      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'VaultController GET /megavault/transfers',
        'Megavault Transfers error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_megavault_transfers.timing`,
        Date.now() - start,
      );
    }
  },
);

router.get(
  '/v1/megavault/transfers/status',
  ...checkSchema({
    transactionHash: {
      in: 'query',
      isString: true,
      notEmpty: true,
      errorMessage: 'transactionHash must be a non-empty string',
    },
  }),
  handleValidationErrors,
  rateLimiterMiddleware(getReqRateLimiter),
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      transactionHash,
    }: MegavaultTransferStatusRequest = matchedData(req) as MegavaultTransferStatusRequest;

    try {
      const controllers: VaultController = new VaultController();
      const response: MegavaultTransferStatusResponse = await controllers
        .getMegavaultTransferStatus(transactionHash);
      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'VaultController GET /megavault/transfers/status',
        'Megavault Transfer Status error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_megavault_transfer_status.timing`,
        Date.now() - start,
      );
    }
  },
);

async function getVaultSubaccountPnlTicks(
  vaultSubaccountIds: string[],
  resolution: PnlTickInterval,
): Promise<PnlTicksFromDatabase[]> {
  if (vaultSubaccountIds.length === 0) {
    return [];
  }

  let windowSeconds: number;
  if (resolution === PnlTickInterval.day) {
    windowSeconds = config.VAULT_PNL_HISTORY_DAYS * 24 * 60 * 60; // days to seconds
  } else {
    windowSeconds = config.VAULT_PNL_HISTORY_HOURS * 60 * 60; // hours to seconds
  }

  const pnlTicks: PnlTicksFromDatabase[] = await VaultPnlTicksView.getVaultsPnl(
    resolution,
    windowSeconds,
    getVaultPnlStartDate(),
  );

  return adjustVaultPnlTicks(pnlTicks, getVaultStartPnl());
}

async function getVaultPositions(
  vaultSubaccounts: VaultMapping,
): Promise<Map<string, VaultPosition>> {
  const start: number = Date.now();
  const vaultSubaccountIds: string[] = _.keys(vaultSubaccounts);
  if (vaultSubaccountIds.length === 0) {
    return new Map();
  }

  const [
    subaccounts,
    assets,
    openPerpetualPositions,
    assetPositions,
    markets,
    latestBlock,
  ]: [
    SubaccountFromDatabase[],
    AssetFromDatabase[],
    PerpetualPositionFromDatabase[],
    AssetPositionFromDatabase[],
    MarketFromDatabase[],
    BlockFromDatabase | undefined,
  ] = await Promise.all([
    SubaccountTable.findAll(
      {
        id: vaultSubaccountIds,
      },
      [],
    ),
    AssetTable.findAll(
      {},
      [],
    ),
    PerpetualPositionTable.findAll(
      {
        subaccountId: vaultSubaccountIds,
        status: [PerpetualPositionStatus.OPEN],
      },
      [],
    ),
    AssetPositionTable.findAll(
      {
        subaccountId: vaultSubaccountIds,
        assetId: [USDC_ASSET_ID],
      },
      [],
    ),
    MarketTable.findAll(
      {},
      [],
    ),
    BlockTable.getLatest(),
  ]);
  stats.timing(
    `${config.SERVICE_NAME}.${controllerName}.positions.fetch_subaccounts_positions.timing`,
    Date.now() - start,
  );

  const startFunding: number = Date.now();
  const updatedAtHeights: string[] = _(subaccounts).map('updatedAtHeight').uniq().value();
  const [
    latestFundingIndexMap,
    fundingIndexMaps,
  ]: [
    FundingIndexMap,
    {[blockHeight: string]: FundingIndexMap}
  ] = await Promise.all([
    FundingIndexUpdatesTable
      .findFundingIndexMap(
        latestBlock.blockHeight,
      ),
    getFundingIndexMapsChunked(updatedAtHeights),
  ]);
  stats.timing(
    `${config.SERVICE_NAME}.${controllerName}.positions.fetch_funding.timing`,
    Date.now() - startFunding,
  );

  const assetPositionsBySubaccount:
  { [subaccountId: string]: AssetPositionFromDatabase[] } = _.groupBy(
    assetPositions,
    'subaccountId',
  );
  const openPerpetualPositionsBySubaccount:
  { [subaccountId: string]: PerpetualPositionFromDatabase[] } = _.groupBy(
    openPerpetualPositions,
    'subaccountId',
  );
  const assetIdToAsset: AssetById = _.keyBy(
    assets,
    AssetColumns.id,
  );

  const vaultPositionsAndSubaccountId: {
    position: VaultPosition,
    subaccountId: string,
  }[] = subaccounts.map((subaccount: SubaccountFromDatabase) => {
    const perpetualMarket: PerpetualMarketFromDatabase | undefined = perpetualMarketRefresher
      .getPerpetualMarketFromClobPairId(vaultSubaccounts[subaccount.id].clobPairId);
    if (perpetualMarket === undefined) {
      throw new Error(
        `Vault clob pair id ${vaultSubaccounts[subaccount.id]} does not correspond to a ` +
          'perpetual market.');
    }
    const lastUpdatedFundingIndexMap: FundingIndexMap = fundingIndexMaps[
      subaccount.updatedAtHeight
    ];
    if (lastUpdatedFundingIndexMap === undefined) {
      throw new Error(
        `No funding indices could be found for vault with subaccount ${subaccount.id}`,
      );
    }

    const subaccountResponse: SubaccountResponseObject = getSubaccountResponse(
      subaccount,
      openPerpetualPositionsBySubaccount[subaccount.id] || [],
      assetPositionsBySubaccount[subaccount.id] || [],
      assets,
      markets,
      perpetualMarketRefresher.getPerpetualMarketsMap(),
      latestBlock.blockHeight,
      latestFundingIndexMap,
      lastUpdatedFundingIndexMap,
    );

    return {
      position: {
        ticker: perpetualMarket.ticker,
        assetPosition: subaccountResponse.assetPositions[
          assetIdToAsset[USDC_ASSET_ID].symbol
        ],
        perpetualPosition: subaccountResponse.openPerpetualPositions[
          perpetualMarket.ticker
        ] || undefined,
        equity: subaccountResponse.equity,
      },
      subaccountId: subaccount.id,
    };
  });

  return new Map(vaultPositionsAndSubaccountId.map(
    (obj: { position: VaultPosition, subaccountId: string }) : [string, VaultPosition] => {
      return [
        obj.subaccountId,
        obj.position,
      ];
    },
  ));
}

async function getMainSubaccountEquity(): Promise<string> {
  // Main vault subaccount should only ever hold a USDC and never any perpetuals.
  const usdcBalance: {[subaccountId: string]: Big} = await AssetPositionTable
    .findUsdcPositionForSubaccounts(
      [MEGAVAULT_SUBACCOUNT_ID],
    );
  return usdcBalance[MEGAVAULT_SUBACCOUNT_ID]?.toFixed() || '0';
}

function getPnlTicksWithCurrentTick(
  equity: string,
  pnlTicks: PnlTicksFromDatabase[],
  latestBlock: BlockFromDatabase,
  latestTick: PnlTicksFromDatabase | undefined = undefined,
): PnlTicksFromDatabase[] {
  if (latestTick !== undefined) {
    return pnlTicks.concat({
      ...latestTick,
      equity,
      blockHeight: latestBlock.blockHeight,
      blockTime: latestBlock.time,
      createdAt: latestBlock.time,
    });
  }
  if (pnlTicks.length === 0) {
    return [];
  }
  const currentTick: PnlTicksFromDatabase = {
    ...(_.maxBy(pnlTicks, 'blockTime')!),
    equity,
    blockHeight: latestBlock.blockHeight,
    blockTime: latestBlock.time,
    createdAt: latestBlock.time,
  };
  return pnlTicks.concat([currentTick]);
}

export async function getLatestPnlTicks(): Promise<PnlTicksFromDatabase[]> {
  const latestPnlTicks: PnlTicksFromDatabase[] = await VaultPnlTicksView.getLatestVaultPnl();
  const adjustedPnlTicks: PnlTicksFromDatabase[] = adjustVaultPnlTicks(
    latestPnlTicks,
    getVaultStartPnl(),
  );
  return adjustedPnlTicks;
}

export async function getLatestPnlTick(
  vaults: VaultFromDatabase[],
): Promise<PnlTicksFromDatabase | undefined> {
  const pnlTicks: PnlTicksFromDatabase[] = await VaultPnlTicksView.getVaultsPnl(
    PnlTickInterval.hour,
    config.VAULT_LATEST_PNL_TICK_WINDOW_HOURS * 60 * 60,
    getVaultPnlStartDate(),
  );
  const adjustedPnlTicks: PnlTicksFromDatabase[] = adjustVaultPnlTicks(
    pnlTicks,
    getVaultStartPnl(),
  );
  // Aggregate and get pnl tick closest to the hour
  const aggregatedTicks: PnlTicksFromDatabase[] = aggregateVaultPnlTicks(
    adjustedPnlTicks,
    vaults,
  );
  const filteredTicks: PnlTicksFromDatabase[] = filterOutIntervalTicks(
    aggregatedTicks,
    PnlTickInterval.hour,
  );
  return _.maxBy(filteredTicks, 'blockTime');
}

/**
 * Takes in an array of PnlTicks and filters out the closest pnl tick per interval.
 * @param pnlTicks Array of pnl ticks.
 * @param resolution Resolution of interval.
 * @returns Array of PnlTicksFromDatabase, one per interval.
 */
function filterOutIntervalTicks(
  pnlTicks: PnlTicksFromDatabase[],
  resolution: PnlTickInterval,
): PnlTicksFromDatabase[] {
  // Track start of intervals to closest Pnl tick.
  const ticksPerInterval: Map<string, PnlTicksFromDatabase> = new Map();
  pnlTicks.forEach((pnlTick: PnlTicksFromDatabase): void => {
    const blockTime: DateTime = DateTime.fromISO(pnlTick.blockTime).toUTC();

    const startOfInterval: DateTime = blockTime.toUTC().startOf(resolution);
    const startOfIntervalStr: string = startOfInterval.toISO();
    const tickForInterval: PnlTicksFromDatabase | undefined = ticksPerInterval.get(
      startOfIntervalStr,
    );
    // No tick for the start of interval, set this tick as the block for the interval.
    if (tickForInterval === undefined) {
      ticksPerInterval.set(startOfIntervalStr, pnlTick);
      return;
    }
    const tickPerIntervalBlockTime: DateTime = DateTime.fromISO(tickForInterval.blockTime);

    // This tick is closer to the start of the interval, set it as the tick for the interval.
    if (blockTime.diff(startOfInterval) < tickPerIntervalBlockTime.diff(startOfInterval)) {
      ticksPerInterval.set(startOfIntervalStr, pnlTick);
    }
  });
  return Array.from(ticksPerInterval.values());
}

function getResolution(resolution: PnlTickInterval = PnlTickInterval.day): PnlTickInterval {
  return resolution;
}

/**
 * Gets funding index maps in a chunked fashion to reduce database load and aggregates into a
 * a map of funding index maps.
 * @param updatedAtHeights
 * @returns
 */
async function getFundingIndexMapsChunked(
  updatedAtHeights: string[],
): Promise<{[blockHeight: string]: FundingIndexMap}> {
  const updatedAtHeightsNum: number[] = updatedAtHeights.map((height: string): number => {
    return parseInt(height, 10);
  }).sort();
  const aggregateFundingIndexMaps: {[blockHeight: string]: FundingIndexMap} = {};
  await Promise.all(getHeightWindows(updatedAtHeightsNum).map(
    async (heightWindow: number[]): Promise<void> => {
      const fundingIndexMaps: {[blockHeight: string]: FundingIndexMap} = await
      FundingIndexUpdatesTable
        .findFundingIndexMaps(
          heightWindow.map((heightNum: number): string => { return heightNum.toString(); }),
        );
      for (const height of _.keys(fundingIndexMaps)) {
        aggregateFundingIndexMaps[height] = fundingIndexMaps[height];
      }
    }));
  return aggregateFundingIndexMaps;
}

/**
 * Separates an array of heights into a chunks based on a window size. Each chunk should only
 * contain heights within a certain number of blocks of each other.
 * @param heights
 * @returns
 */
function getHeightWindows(
  heights: number[],
): number[][] {
  if (heights.length === 0) {
    return [];
  }
  const windows: number[][] = [];
  let windowStart: number = heights[0];
  let currentWindow: number[] = [];
  for (const height of heights) {
    if (height - windowStart < config.VAULT_FETCH_FUNDING_INDEX_BLOCK_WINDOWS) {
      currentWindow.push(height);
    } else {
      windows.push(currentWindow);
      currentWindow = [height];
      windowStart = height;
    }
  }
  windows.push(currentWindow);
  return windows;
}

async function getFirstMainVaultTransferDateTime(): Promise<DateTime | undefined> {
  const { results }: {
    results: TransferFromDatabase[],
  } = await TransferTable.findAllToOrFromSubaccountId(
    {
      subaccountId: [MEGAVAULT_SUBACCOUNT_ID],
      limit: 1,
    },
    [],
    {
      orderBy: [[TransferColumns.createdAt, Ordering.ASC]],
    },
  );
  if (results.length === 0) {
    return undefined;
  }
  return DateTime.fromISO(results[0].createdAt);
}

/**
 * Aggregates vault pnl ticks per hour, filtering out pnl ticks made up of less ticks than expected.
 * Expected number of pnl ticks is calculated from the number of vaults that were created before
 * the pnl tick was created.
 * @param vaultPnlTicks Pnl ticks to aggregate.
 * @param vaults List of all valid vaults.
 * @param mainVaultCreatedAt Date time when the main vault was created or undefined if it does not
 * exist yet.
 * @returns
 */
function aggregateVaultPnlTicks(
  vaultPnlTicks: PnlTicksFromDatabase[],
  vaults: VaultFromDatabase[],
  mainVaultCreatedAt?: DateTime,
): PnlTicksFromDatabase[] {
  // aggregate pnlTicks for all vault subaccounts grouped by blockHeight
  const aggregatedPnlTicks: AggregatedPnlTick[] = aggregateHourlyPnlTicks(vaultPnlTicks);
  const vaultCreationTimes: DateTime[] = _.map(vaults, 'createdAt').map(
    (createdAt: string) => { return DateTime.fromISO(createdAt); },
  ).concat(
    mainVaultCreatedAt === undefined ? [] : [mainVaultCreatedAt],
  ).sort(
    (a: DateTime, b: DateTime) => {
      return a.diff(b).milliseconds;
    },
  );
  return aggregatedPnlTicks.filter((aggregatedTick: AggregatedPnlTick) => {
    // Get number of vaults created before the pnl tick was created by binary-searching for the
    // index of the pnl ticks createdAt in a sorted array of vault createdAt times.
    const numVaultsCreated: number = bounds.le(
      vaultCreationTimes,
      DateTime.fromISO(aggregatedTick.pnlTick.createdAt),
      (a: DateTime, b: DateTime) => { return a.diff(b).milliseconds; },
    );
    // Number of ticks should be greater than number of vaults created before it as there should be
    // a tick for the main vault subaccount.
    return aggregatedTick.numTicks >= numVaultsCreated;
  }).map((aggregatedPnlTick: AggregatedPnlTick) => { return aggregatedPnlTick.pnlTick; });
}

function adjustVaultPnlTicks(
  pnlTicks: PnlTicksFromDatabase[],
  pnlTicksToAdjustBy: PnlTicksFromDatabase[],
): PnlTicksFromDatabase[] {
  const subaccountToPnlTick: {[subaccountId: string]: PnlTicksFromDatabase} = {};
  for (const pnlTickToAdjustBy of pnlTicksToAdjustBy) {
    subaccountToPnlTick[pnlTickToAdjustBy.subaccountId] = pnlTickToAdjustBy;
  }

  return pnlTicks.map((pnlTick: PnlTicksFromDatabase): PnlTicksFromDatabase => {
    const adjustByPnlTick: PnlTicksFromDatabase | undefined = subaccountToPnlTick[
      pnlTick.subaccountId
    ];
    if (adjustByPnlTick === undefined) {
      return pnlTick;
    }
    return {
      ...pnlTick,
      totalPnl: Big(pnlTick.totalPnl).sub(Big(adjustByPnlTick.totalPnl)).toFixed(),
    };
  });
}

/**
 * Computes the 30-day annualized return of the megavault from an ascending series of
 * aggregated PnL ticks. The baseline is the latest tick at least 30 days older than the
 * newest tick, falling back to the oldest available tick. Returns null when there is less
 * than a full day of history or the baseline equity is non-positive.
 */
function computeAnnualizedReturn(
  pnlTicks: PnlTicksResponseObject[],
): string | null {
  const latestTick: PnlTicksResponseObject | undefined = _.last(pnlTicks);
  if (latestTick === undefined) {
    return null;
  }

  const latestTime: DateTime = DateTime.fromISO(latestTick.blockTime).toUTC();
  const aprWindowStart: DateTime = latestTime.minus({ days: 30 });
  const baselineTick: PnlTicksResponseObject = _.findLast(
    pnlTicks,
    (pnlTick: PnlTicksResponseObject): boolean => {
      return DateTime.fromISO(pnlTick.blockTime).toUTC() <= aprWindowStart;
    },
  ) ?? pnlTicks[0];

  const elapsedSeconds: number = latestTime.diff(
    DateTime.fromISO(baselineTick.blockTime).toUTC(),
  ).as('seconds');
  const baselineEquity: Big = Big(baselineTick.equity);
  if (elapsedSeconds < 24 * 60 * 60 || baselineEquity.lte(0)) {
    return null;
  }

  return Big(latestTick.totalPnl)
    .sub(baselineTick.totalPnl)
    .div(baselineEquity)
    .times(365 * 24 * 60 * 60)
    .div(elapsedSeconds)
    .toFixed();
}

/**
 * Computes the largest peak-to-trough decline of cumulative PnL (in USDC) over an ascending
 * series of aggregated PnL ticks. Computed on PnL rather than equity so that deposits and
 * withdrawals do not register as gains or drawdowns.
 */
function computeMaxPnlDrawdown(
  pnlTicks: PnlTicksResponseObject[],
): string {
  let maxDrawdown: Big = Big(0);
  let peakPnl: Big | undefined;

  for (const pnlTick of pnlTicks) {
    const pnl: Big = Big(pnlTick.totalPnl);
    if (peakPnl === undefined || pnl.gt(peakPnl)) {
      peakPnl = pnl;
    }
    const drawdown: Big = peakPnl.minus(pnl);
    if (drawdown.gt(maxDrawdown)) {
      maxDrawdown = drawdown;
    }
  }

  return maxDrawdown.toFixed();
}

function megavaultTransferToResponseObject(
  transfer: TransferFromDatabase,
  assetMap: AssetById,
  subaccountsById: { [id: string]: SubaccountFromDatabase },
): MegavaultTransferResponseObject {
  const isDeposit: boolean = transfer.recipientSubaccountId === MEGAVAULT_SUBACCOUNT_ID;
  const userSubaccountId: string | undefined = isDeposit
    ? transfer.senderSubaccountId
    : transfer.recipientSubaccountId;
  const userWalletAddress: string | undefined = isDeposit
    ? transfer.senderWalletAddress
    : transfer.recipientWalletAddress;
  const userSubaccount: SubaccountFromDatabase | undefined = userSubaccountId === undefined
    ? undefined
    : subaccountsById[userSubaccountId];

  return {
    id: transfer.id,
    type: isDeposit ? MegavaultTransferType.DEPOSIT : MegavaultTransferType.WITHDRAWAL,
    address: userWalletAddress ?? userSubaccount?.address ?? '',
    subaccountNumber: userSubaccount?.subaccountNumber,
    size: transfer.size,
    symbol: assetMap[transfer.assetId].symbol,
    createdAt: transfer.createdAt,
    createdAtHeight: transfer.createdAtHeight,
    transactionHash: transfer.transactionHash,
  };
}

export default router;
