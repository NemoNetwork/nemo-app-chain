/**
 * Environment variables required by Ender.
 */

import {
  parseSchema,
  baseConfigSchema,
  parseBoolean,
} from '@nemo-network-indexer/base/src';
import {
  kafkaConfigSchema,
} from '@nemo-network-indexer/kafka/src';
import {
  postgresConfigSchema,
} from '@nemo-network-indexer/postgres/src';
import { redisConfigSchema } from '@nemo-network-indexer/redis/src';

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
