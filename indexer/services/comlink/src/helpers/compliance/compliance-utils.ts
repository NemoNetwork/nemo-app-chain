import { CountryHeaders, isRestrictedCountryHeaders } from '@nemo-network-indexer/compliance/src';
import { ComplianceReason } from '@nemo-network-indexer/postgres/build/src';

export function getGeoComplianceReason(
  headers: CountryHeaders,
): ComplianceReason | undefined {
  if (isRestrictedCountryHeaders(headers)) {
    const country: string | undefined = headers['cf-ipcountry'];
    if (country === 'US') {
      return ComplianceReason.US_GEO;
    } else if (country === 'CA') {
      return ComplianceReason.CA_GEO;
    } else if (country === 'GB') {
      return ComplianceReason.GB_GEO;
    } else {
      return ComplianceReason.SANCTIONED_GEO;
    }
  }
  return undefined;
}
