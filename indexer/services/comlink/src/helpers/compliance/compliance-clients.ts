import {
  ComplianceClient,
  getComplianceClient,
} from '@nemo-network-indexer/compliance';
import { ComplianceProvider } from '@nemo-network-indexer/postgres';

export interface ClientAndProvider {
  client: ComplianceClient;
  provider: ComplianceProvider;
}

export const complianceProvider: ClientAndProvider = {
  client: getComplianceClient(),
  provider: ComplianceProvider.ELLIPTIC,
};
