/**
 * Environment variables required by Ender.
 */

import {
  parseSchema,
  baseConfigSchema,
  parseBoolean,
} from '@nemo-network-indexer/base/build';
import {
  kafkaConfigSchema,
} from '@nemo-network-indexer/kafka/build/src';
import {
  postgresConfigSchema,
} from '@nemo-network-indexer/postgres/build/src';
import { redisConfigSchema } from '@nemo-network-indexer/redis/build/redis/src';

export const configSchema = {
  ...baseConfigSchema,
  ...postgresConfigSchema,
  ...redisConfigSchema,
  ...kafkaConfigSchema,
  SEND_WEBSOCKET_MESSAGES: parseBoolean({
    default: true,
  }),
};

export default parseSchema(configSchema);
