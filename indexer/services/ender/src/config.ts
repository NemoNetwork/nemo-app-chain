/**
 * Environment variables required by Ender.
 */

import {
  parseSchema,
  baseConfigSchema,
  parseBoolean,
} from '@nemo-network-indexer/base';
import {
  kafkaConfigSchema,
} from '@nemo-network-indexer/kafka';
import {
  postgresConfigSchema,
} from '@nemo-network-indexer/postgres';
import { redisConfigSchema } from '@nemo-network-indexer/redis';

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
