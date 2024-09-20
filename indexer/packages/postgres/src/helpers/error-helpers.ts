import { logger } from '@nemo-network-indexer/base/build';

import { ValidationError } from '../lib/errors';

export function logAndThrowValidationError(message: string): void {
  logger.error({
    at: 'errorHelpers#logAndThrowValidationError',
    message,
  });
  throw new ValidationError(message);
}
