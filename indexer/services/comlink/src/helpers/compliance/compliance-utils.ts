import { CountryHeaders, isRestrictedCountryHeaders } from '@nemo-network-indexer/compliance';
import { ComplianceReason } from '@nemo-network-indexer/postgres';

export function getGeoComplianceReason(
  headers: CountryHeaders,
): ComplianceReason | undefined {
  if (isRestrictedCountryHeaders(headers)) {
    const country: string | undefined = headers['cf-ipcountry'];
    if (country === 'US') {
      return ComplianceReason.US_GEO;
    } else if (country === 'CA') {
      return ComplianceReason.CA_GEO;
    } else {
      return ComplianceReason.SANCTIONED_GEO;
    }
  }
  return undefined;
}
