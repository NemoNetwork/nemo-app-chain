import { testConstants } from '@nemo-network-indexer/postgres/src';
import { PnlTickForSubaccounts } from '@nemo-network-indexer/redis/src';

export const defaultPnlTickForSubaccounts: PnlTickForSubaccounts = {
  [testConstants.defaultSubaccountId]: testConstants.defaultPnlTick,
  [testConstants.defaultSubaccountId2]: {
    ...testConstants.defaultPnlTick,
    subaccountId: testConstants.defaultSubaccountId2,
    equity: '9000',
  },
};
