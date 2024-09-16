import {
  ComplianceClient,
  getComplianceClient,
} from '@nemo-network-indexer/compliance/src';
import { ComplianceProvider } from '@nemo-network-indexer/postgres/src';

export interface ClientAndProvider {
  client: ComplianceClient,
  provider: ComplianceProvider,
}

export const complianceProvider: ClientAndProvider = {
  client: getComplianceClient(),
  provider: ComplianceProvider.ELLIPTIC,
};
