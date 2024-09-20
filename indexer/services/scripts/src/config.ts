import { baseConfigSchema, parseSchema } from '@nemo-network-indexer/base/build';
import { kafkaConfigSchema } from '@nemo-network-indexer/kafka/build/src';
import { postgresConfigSchema } from '@nemo-network-indexer/postgres/build/src';

export const configSchema = {
  ...baseConfigSchema,
  ...postgresConfigSchema,
  ...kafkaConfigSchema,
};

export default parseSchema(configSchema);
