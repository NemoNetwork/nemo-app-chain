import {
  dbHelpers,
  testConstants,
  testMocks,
  perpetualMarketRefresher,
  liquidityTierRefresher,
  BlockTable,
  SubaccountTable,
  AssetPositionTable,
  PerpetualPositionTable,
  FundingIndexUpdatesTable,
  PnlTicksTable,
  FillTable,
  TransferTable,
  VaultTable,
  VaultPnlTicksView,
  MEGAVAULT_MODULE_ADDRESS,
  MEGAVAULT_SUBACCOUNT_ID,
} from '@nemo-network-indexer/postgres';
import request from 'supertest';
import { DateTime, Settings } from 'luxon';
import Big from 'big.js';

import { RequestMethod } from '../../../../src/types';
import { sendRequest } from '../../../helpers/helpers';
import config from '../../../../src/config';
import { clearVaultStartPnl, startVaultStartPnlCache } from '../../../../src/caches/vault-start-pnl';

describe('vault-controller#V4 megavault endpoints', () => {
  const latestBlockHeight: string = '25';
  const currentDayBlockHeight: string = '9';
  const twoHourBlockHeight: string = '7';
  const twoDayBlockHeight: string = '3';
  const currentDay: DateTime = DateTime.utc().startOf('day').minus({ hour: 5 });
  const latestTime: DateTime = currentDay.plus({ minute: 90 });
  const twoHoursAgo: DateTime = currentDay.minus({ hour: 2 });
  const twoDaysAgo: DateTime = currentDay.minus({ day: 2 });
  const initialFundingIndex: string = '10000';
  const vault1Equity: number = 159500;
  const mainVaultEquity: number = 10000;
  const vaultPnlHistoryHoursPrev: number = config.VAULT_PNL_HISTORY_HOURS;
  const vaultPnlLastPnlWindowPrev: number = config.VAULT_LATEST_PNL_TICK_WINDOW_HOURS;
  const vaultPnlStartDatePrev: string = config.VAULT_PNL_START_DATE;

  const depositTxHash: string = 'deposittxhash';
  const withdrawalTxHash: string = 'withdrawaltxhash';
  const unrelatedTxHash: string = 'unrelatedtxhash';

  beforeAll(async () => {
    await dbHelpers.migrate();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
  });

  beforeEach(async () => {
    config.VAULT_PNL_HISTORY_HOURS = 168;
    config.VAULT_LATEST_PNL_TICK_WINDOW_HOURS = 48;
    config.VAULT_PNL_START_DATE = '2020-01-01T00:00:00Z';
    await testMocks.seedData();
    await perpetualMarketRefresher.updatePerpetualMarkets();
    await liquidityTierRefresher.updateLiquidityTiers();
    await Promise.all([
      BlockTable.create({
        ...testConstants.defaultBlock,
        time: twoDaysAgo.toISO(),
        blockHeight: twoDayBlockHeight,
      }),
      BlockTable.create({
        ...testConstants.defaultBlock,
        time: twoHoursAgo.toISO(),
        blockHeight: twoHourBlockHeight,
      }),
      BlockTable.create({
        ...testConstants.defaultBlock,
        time: currentDay.toISO(),
        blockHeight: currentDayBlockHeight,
      }),
      BlockTable.create({
        ...testConstants.defaultBlock,
        time: latestTime.toISO(),
        blockHeight: latestBlockHeight,
      }),
    ]);
    await SubaccountTable.create({
      address: MEGAVAULT_MODULE_ADDRESS,
      subaccountNumber: 0,
      updatedAt: latestTime.toISO(),
      updatedAtHeight: latestBlockHeight,
    });
    await Promise.all([
      PerpetualPositionTable.create(
        testConstants.defaultPerpetualPosition,
      ),
      AssetPositionTable.upsert(testConstants.defaultAssetPosition),
      FundingIndexUpdatesTable.create({
        ...testConstants.defaultFundingIndexUpdate,
        fundingIndex: initialFundingIndex,
        effectiveAtHeight: testConstants.createdHeight,
      }),
      FundingIndexUpdatesTable.create({
        ...testConstants.defaultFundingIndexUpdate,
        eventId: testConstants.defaultTendermintEventId2,
        effectiveAtHeight: twoDayBlockHeight,
      }),
    ]);
    Settings.now = () => latestTime.valueOf();
  });

  afterEach(async () => {
    await dbHelpers.clearData();
    await VaultPnlTicksView.refreshDailyView();
    await VaultPnlTicksView.refreshHourlyView();
    clearVaultStartPnl();
    config.VAULT_PNL_HISTORY_HOURS = vaultPnlHistoryHoursPrev;
    config.VAULT_LATEST_PNL_TICK_WINDOW_HOURS = vaultPnlLastPnlWindowPrev;
    config.VAULT_PNL_START_DATE = vaultPnlStartDatePrev;
    Settings.now = () => new Date().valueOf();
  });

  describe('GET /vaults', () => {
    it('returns an empty list with no vaults', async () => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: '/vault/v1/vaults',
      });

      expect(response.body).toEqual({ vaults: [] });
    });

    it('returns vaults with tickers and skips vaults with unknown clob pairs', async () => {
      await Promise.all([
        VaultTable.create({
          ...testConstants.defaultVault,
          address: testConstants.defaultSubaccount.address,
          clobPairId: testConstants.defaultPerpetualMarket.clobPairId,
        }),
        VaultTable.create({
          ...testConstants.defaultVault,
          address: testConstants.vaultAddress,
          clobPairId: '999',
        }),
      ]);

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: '/vault/v1/vaults',
      });

      expect(response.body.vaults).toHaveLength(1);
      expect(response.body.vaults[0]).toEqual({
        address: testConstants.defaultSubaccount.address,
        ticker: testConstants.defaultPerpetualMarket.ticker,
        status: testConstants.defaultVault.status,
        createdAt: testConstants.defaultVault.createdAt,
        updatedAt: testConstants.defaultVault.updatedAt,
      });
    });
  });

  describe('GET /megavault/summary', () => {
    it('returns zeroed summary with no vaults or activity', async () => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: '/vault/v1/megavault/summary',
      });

      expect(response.body).toEqual({
        equity: '0',
        numVaults: 0,
        allTimePnl: '0',
        apr: null,
        maxDrawdown: '0',
        volume24H: '0',
        createdAt: null,
      });
    });

    it('returns equity, pnl metrics, volume and creation time with a vault', async () => {
      await VaultTable.create({
        ...testConstants.defaultVault,
        address: testConstants.defaultSubaccount.address,
        clobPairId: testConstants.defaultPerpetualMarket.clobPairId,
      });
      await AssetPositionTable.upsert({
        ...testConstants.defaultAssetPosition,
        subaccountId: MEGAVAULT_SUBACCOUNT_ID,
      });
      // PnL goes 10000 -> 2000: all-time PnL is 2000 and max drawdown is 8000.
      await Promise.all([
        PnlTicksTable.create({
          ...testConstants.defaultPnlTick,
          equity: '20000',
          totalPnl: '10000',
          blockTime: twoDaysAgo.toISO(),
          createdAt: twoDaysAgo.toISO(),
          blockHeight: twoDayBlockHeight,
        }),
        PnlTicksTable.create({
          ...testConstants.defaultPnlTick,
          equity: '12000',
          totalPnl: '2000',
          blockTime: currentDay.toISO(),
          createdAt: currentDay.toISO(),
          blockHeight: currentDayBlockHeight,
        }),
      ]);
      await VaultPnlTicksView.refreshDailyView();
      await VaultPnlTicksView.refreshHourlyView();
      await startVaultStartPnlCache();
      await Promise.all([
        // Fill within the past 24 hours: counts towards volume.
        FillTable.create({
          ...testConstants.defaultFill,
          orderId: undefined,
          price: '20000',
          size: '0.5',
          eventId: testConstants.defaultTendermintEventId3,
          createdAt: twoHoursAgo.toISO(),
          createdAtHeight: twoHourBlockHeight,
        }),
        // Fill older than 24 hours: excluded from volume.
        FillTable.create({
          ...testConstants.defaultFill,
          orderId: undefined,
          price: '20000',
          size: '1',
          eventId: testConstants.defaultTendermintEventId4,
          createdAt: twoDaysAgo.toISO(),
          createdAtHeight: twoDayBlockHeight,
        }),
        TransferTable.create({
          ...testConstants.defaultTransfer,
          senderSubaccountId: testConstants.defaultSubaccountId,
          recipientSubaccountId: MEGAVAULT_SUBACCOUNT_ID,
          transactionHash: depositTxHash,
          createdAt: twoDaysAgo.toISO(),
          createdAtHeight: twoDayBlockHeight,
        }),
      ]);

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: '/vault/v1/megavault/summary',
      });

      const elapsedSeconds: number = latestTime.diff(twoDaysAgo).as('seconds');
      const expectedApr: string = Big('2000')
        .sub('10000')
        .div('20000')
        .times(365 * 24 * 60 * 60)
        .div(elapsedSeconds)
        .toFixed();

      expect(response.body).toEqual(expect.objectContaining({
        equity: Big(vault1Equity).add(mainVaultEquity).toFixed(),
        numVaults: 1,
        allTimePnl: '2000',
        apr: expectedApr,
        maxDrawdown: '8000',
        volume24H: '10000',
      }));
      expect(DateTime.fromISO(response.body.createdAt).toMillis()).toEqual(twoDaysAgo.toMillis());
    });
  });

  describe('GET /megavault/transfers', () => {
    async function createTransfers(): Promise<void> {
      const alternateSubaccountId: string = await SubaccountTable.create(
        testConstants.defaultSubaccountWithAlternateAddress,
      ).then((subaccount) => subaccount.id);
      await Promise.all([
        // Deposit: user subaccount -> megavault.
        TransferTable.create({
          ...testConstants.defaultTransfer,
          senderSubaccountId: testConstants.defaultSubaccountId,
          recipientSubaccountId: MEGAVAULT_SUBACCOUNT_ID,
          size: '100',
          transactionHash: depositTxHash,
          createdAt: twoDaysAgo.toISO(),
          createdAtHeight: twoDayBlockHeight,
        }),
        // Withdrawal: megavault -> user subaccount.
        TransferTable.create({
          ...testConstants.defaultTransfer,
          senderSubaccountId: MEGAVAULT_SUBACCOUNT_ID,
          recipientSubaccountId: testConstants.defaultSubaccountId,
          size: '50',
          eventId: testConstants.defaultTendermintEventId2,
          transactionHash: withdrawalTxHash,
          createdAt: twoHoursAgo.toISO(),
          createdAtHeight: twoHourBlockHeight,
        }),
        // Unrelated transfer between user subaccounts: never included.
        TransferTable.create({
          ...testConstants.defaultTransfer,
          size: '25',
          eventId: testConstants.defaultTendermintEventId3,
          transactionHash: unrelatedTxHash,
          createdAt: twoHoursAgo.toISO(),
          createdAtHeight: twoHourBlockHeight,
        }),
        // Deposit from a different user: excluded when filtering by address.
        TransferTable.create({
          ...testConstants.defaultTransfer,
          senderSubaccountId: alternateSubaccountId,
          recipientSubaccountId: MEGAVAULT_SUBACCOUNT_ID,
          size: '75',
          eventId: testConstants.defaultTendermintEventId4,
          transactionHash: 'otherusertxhash',
          createdAt: twoHoursAgo.toISO(),
          createdAtHeight: twoHourBlockHeight,
        }),
      ]);
    }

    it('returns deposits and withdrawals for an address, newest first', async () => {
      await createTransfers();

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/vault/v1/megavault/transfers?address=${testConstants.defaultAddress}`,
      });

      expect(response.body.transfers).toHaveLength(2);
      expect(response.body.transfers[0]).toEqual(expect.objectContaining({
        type: 'WITHDRAWAL',
        address: testConstants.defaultAddress,
        subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
        size: '50',
        symbol: testConstants.defaultAsset.symbol,
        createdAt: twoHoursAgo.toISO(),
        createdAtHeight: twoHourBlockHeight,
        transactionHash: withdrawalTxHash,
      }));
      expect(response.body.transfers[1]).toEqual(expect.objectContaining({
        type: 'DEPOSIT',
        address: testConstants.defaultAddress,
        subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
        size: '100',
        symbol: testConstants.defaultAsset.symbol,
        createdAt: twoDaysAgo.toISO(),
        createdAtHeight: twoDayBlockHeight,
        transactionHash: depositTxHash,
      }));
    });

    it('respects the limit parameter', async () => {
      await createTransfers();

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/vault/v1/megavault/transfers?address=${testConstants.defaultAddress}&limit=1`,
      });

      expect(response.body.transfers).toHaveLength(1);
      expect(response.body.transfers[0].transactionHash).toEqual(withdrawalTxHash);
    });

    it('returns an empty list for an address with no subaccounts', async () => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: '/vault/v1/megavault/transfers?address=nemounknownaddress',
      });

      expect(response.body).toEqual({ transfers: [] });
    });

    it('returns 400 when address is missing', async () => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: '/vault/v1/megavault/transfers',
        expectedStatus: 400,
      });

      expect(response.status).toEqual(400);
    });
  });

  describe('GET /megavault/transfers/status', () => {
    it('returns COMPLETED with the transfer when the hash is indexed', async () => {
      await TransferTable.create({
        ...testConstants.defaultTransfer,
        senderSubaccountId: testConstants.defaultSubaccountId,
        recipientSubaccountId: MEGAVAULT_SUBACCOUNT_ID,
        size: '100',
        transactionHash: depositTxHash,
        createdAt: twoDaysAgo.toISO(),
        createdAtHeight: twoDayBlockHeight,
      });

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/vault/v1/megavault/transfers/status?transactionHash=${depositTxHash}`,
      });

      expect(response.body.status).toEqual('COMPLETED');
      expect(response.body.transactionHash).toEqual(depositTxHash);
      expect(response.body.transfers).toHaveLength(1);
      expect(response.body.transfers[0]).toEqual(expect.objectContaining({
        type: 'DEPOSIT',
        address: testConstants.defaultAddress,
        size: '100',
        transactionHash: depositTxHash,
      }));
    });

    it('returns PENDING for an unknown hash', async () => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: '/vault/v1/megavault/transfers/status?transactionHash=doesnotexist',
      });

      expect(response.body).toEqual({
        transactionHash: 'doesnotexist',
        status: 'PENDING',
        transfers: [],
      });
    });

    it('returns PENDING for a hash whose transfers do not involve the megavault', async () => {
      await TransferTable.create({
        ...testConstants.defaultTransfer,
        size: '25',
        transactionHash: unrelatedTxHash,
        createdAt: twoHoursAgo.toISO(),
        createdAtHeight: twoHourBlockHeight,
      });

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/vault/v1/megavault/transfers/status?transactionHash=${unrelatedTxHash}`,
      });

      expect(response.body.status).toEqual('PENDING');
      expect(response.body.transfers).toEqual([]);
    });
  });
});
