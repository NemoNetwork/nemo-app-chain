import { logger } from '@nemo-network-indexer/base/build';
import {
  WebsocketTopics,
  consumer,
  stopConsumer,
} from '@nemo-network-indexer/kafka/build/src';

export async function connect(): Promise<void> {
  await consumer.connect();

  logger.info({
    at: 'kafka-controller#connect',
    message: 'Connected to Kafka',
  });

  await consumer.subscribe({ topics: Object.values(WebsocketTopics) });
}

export async function disconnect(): Promise<void> {
  logger.info({
    at: 'kafka-controller#stop',
    message: 'Stopping kafka consumer',
  });

  await stopConsumer();
}
