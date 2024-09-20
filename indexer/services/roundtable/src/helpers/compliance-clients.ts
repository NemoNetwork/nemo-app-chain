import {
  ComplianceClient,
  getComplianceClient,
} from '@nemo-network-indexer/compliance/build';
import { ComplianceProvider } from '@nemo-network-indexer/postgres/build/src';

export interface ClientAndProvider {
  client: ComplianceClient,
  provider: ComplianceProvider,
}

export const complianceProvider: ClientAndProvider = {
  client: getComplianceClient(),
  provider: ComplianceProvider.ELLIPTIC,
};
