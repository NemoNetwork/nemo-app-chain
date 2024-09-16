import { baseConfigSchema, parseSchema } from '@nemo-network-indexer/base/src';
import { kafkaConfigSchema } from '@nemo-network-indexer/kafka/src';
import { postgresConfigSchema } from '@nemo-network-indexer/postgres/src';

export const configSchema = {
  ...baseConfigSchema,
  ...postgresConfigSchema,
  ...kafkaConfigSchema,
};

export default parseSchema(configSchema);
