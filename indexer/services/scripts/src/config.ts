import { baseConfigSchema, parseSchema } from '@nemo-network-indexer/base';
import { kafkaConfigSchema } from '@nemo-network-indexer/kafka';
import { postgresConfigSchema } from '@nemo-network-indexer/postgres';

export const configSchema = {
  ...baseConfigSchema,
  ...postgresConfigSchema,
  ...kafkaConfigSchema,
};

export default parseSchema(configSchema);
