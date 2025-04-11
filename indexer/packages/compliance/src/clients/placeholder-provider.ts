import { ComplianceClientResponse } from '../types';
import { ComplianceClient } from './compliance-client';

export class PlaceHolderProviderClient extends ComplianceClient {
  public getComplianceResponse(address: string): Promise<ComplianceClientResponse> {
    return Promise.resolve({
      address,
      blocked: false,
      riskScore: '0',
    });
  }
}
