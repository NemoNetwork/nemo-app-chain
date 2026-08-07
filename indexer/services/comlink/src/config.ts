import {
  baseConfigSchema,
  parseBoolean,
  parseInteger,
  parseSchema,
  parseString,
} from '@nemo-network-indexer/base/build';
import { complianceConfigSchema } from '@nemo-network-indexer/compliance/build';
import {
  postgresConfigSchema,
} from '@nemo-network-indexer/postgres/build/src';
import { redisConfigSchema } from '@nemo-network-indexer/redis/build/redis/src';

export const configSchema = {
  ...baseConfigSchema,
  ...postgresConfigSchema,
  ...redisConfigSchema,
  ...complianceConfigSchema,

  CHAIN_ID: parseString({ default: 'nemo-network' }),
  API_LIMIT_V4: parseInteger({
    default: 1000,
  }),
  API_ORDERBOOK_LEVELS_PER_SIDE_LIMIT: parseInteger({ default: 100 }),

  // Logging config
  LOG_GETS: parseBoolean({ default: false }),

  // Express server config
  PORT: parseInteger({ default: 8080 }),
  CORS_ORIGIN: parseString({ default: '*' }),
  KEEP_ALIVE_MS: parseInteger({ default: 61_000 }),
  HEADERS_TIMEOUT_MS: parseInteger({ default: 65_000 }),

  // Rate limit Redis URL
  RATE_LIMIT_REDIS_URL: parseString({
    default: 'redis://redis:6379',
  }),
  // Rate limits
  RATE_LIMIT_ENABLED: parseBoolean({ default: true }),
  // IP addresses internal to the Indexer have no rate-limit
  INDEXER_INTERNAL_IPS: parseString({ default: '' }),
  // Points / duration determines the maximum rate of requests given that each requests costs 1
  // point
  RATE_LIMIT_GET_POINTS: parseInteger({ default: 100 }),
  RATE_LIMIT_GET_DURATION_SECONDS: parseInteger({ default: 10 }), // 100 requests / 10 seconds

  // Rate limit for screening new / refreshed addresses
  RATE_LIMIT_SCREEN_QUERY_PROVIDER_POINTS: parseInteger({ default: 2 }),
  RATE_LIMIT_SCREEN_QUERY_PROVIDER_DURATION_SECONDS: parseInteger({ default: 60 }), // 2 reqs / min
  RATE_LIMIT_SCREEN_QUERY_PROVIDER_GLOBAL_POINTS: parseInteger({ default: 100 }),
  // 100 req / min
  RATE_LIMIT_SCREEN_QUERY_PROVIDER_GLOBAL_DURATION_SECONDS: parseInteger({ default: 60 }),
  // Threshold for refreshing compliance data for an address when screened
  MAX_AGE_SCREENED_ADDRESS_COMPLIANCE_DATA_SECONDS: parseInteger({ default: 86_400 }), //  1 day
  // Expose setting compliance status, only set to true in dev/staging.
  EXPOSE_SET_COMPLIANCE_ENDPOINT: parseBoolean({ default: false }),

  // Vaults config
  VAULT_PNL_HISTORY_DAYS: parseInteger({ default: 90 }),
  VAULT_PNL_HISTORY_HOURS: parseInteger({ default: 72 }),
  // Vault pnl before this date is not served. It exists so that vaults created at different
  // times can be rebased onto a common start and summed without the oldest dominating.
  VAULT_PNL_START_DATE: parseString({ default: '2024-01-01T00:00:00Z' }),
  VAULT_LATEST_PNL_TICK_WINDOW_HOURS: parseInteger({ default: 1 }),
  VAULT_FETCH_FUNDING_INDEX_BLOCK_WINDOWS: parseInteger({ default: 250_000 }),

};

////////////////////////////////////////////////////////////////////////////////
//                             CONFIG PROCESSING                              //
////////////////////////////////////////////////////////////////////////////////

// Process the top-level configuration.
const config = parseSchema(configSchema);

export default config;
