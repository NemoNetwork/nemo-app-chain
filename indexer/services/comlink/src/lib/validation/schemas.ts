import {
  perpetualMarketRefresher,
  MAX_PARENT_SUBACCOUNTS,
  CHILD_SUBACCOUNT_MULTIPLIER,
} from '@nemo-network-indexer/postgres/build/src';
import { checkSchema, ParamSchema } from 'express-validator';
import { DateTime } from 'luxon';

import config from '../../config';

/**
 * Sanitizes a date string to ISO 8601 format.
 * Accepts both ISO 8601 format and U.S. datetime formats (MM/DD/YYYY, MM-DD-YYYY, etc.)
 * @param value - The date string to sanitize
 * @returns ISO 8601 formatted date string, or the original value if parsing fails
 */
function sanitizeDateToISO(value?: string): string | undefined {
  if (!value || typeof value !== 'string') {
    return value;
  }

  // First, try to parse as ISO 8601 (current format)
  let dt = DateTime.fromISO(value);
  if (dt.isValid) {
    return dt.toISO();
  }

  // Try U.S. datetime formats
  // MM/DD/YYYY or MM/DD/YYYY HH:mm:ss or MM/DD/YYYY HH:mm:ss AM/PM
  const usFormats = [
    'MM/dd/yyyy',
    'MM/dd/yyyy HH:mm:ss',
    'MM/dd/yyyy hh:mm:ss a',
    'MM/dd/yyyy HH:mm',
    'MM/dd/yyyy hh:mm a',
    'MM-dd-yyyy',
    'MM-dd-yyyy HH:mm:ss',
    'MM-dd-yyyy hh:mm:ss a',
    'MM-dd-yyyy HH:mm',
    'MM-dd-yyyy hh:mm a',
    'M/d/yyyy',
    'M/d/yyyy HH:mm:ss',
    'M/d/yyyy hh:mm:ss a',
    'M/d/yyyy HH:mm',
    'M/d/yyyy hh:mm a',
  ];

  for (const format of usFormats) {
    dt = DateTime.fromFormat(value, format);
    if (dt.isValid) {
      return dt.toISO();
    }
  }

  // If all parsing fails, return original value (will be caught by isISO8601 validator)
  return value;
}

export const CheckSubaccountSchema = checkSchema({
  address: {
    in: ['params', 'query'],
    isString: true,
  },
  subaccountNumber: {
    in: ['params', 'query'],
    isInt: {
      options: { gt: -1, lt: MAX_PARENT_SUBACCOUNTS * CHILD_SUBACCOUNT_MULTIPLIER + 1 },
    },
    errorMessage: 'subaccountNumber must be a non-negative integer less than 128001',
  },
});

export const CheckParentSubaccountSchema = checkSchema({
  address: {
    in: ['params', 'query'],
    isString: true,
  },
  parentSubaccountNumber: {
    in: ['params', 'query'],
    isInt: {
      options: { gt: -1, lt: MAX_PARENT_SUBACCOUNTS },
    },
    errorMessage: 'parentSubaccountNumber must be a non-negative integer less than 128',
  },
});

export const checkAddressSchemaRecord: Record<string, ParamSchema> = {
  address: {
    in: ['params'],
    isString: true,
  },
};

export const CheckAddressSchema = checkSchema(checkAddressSchemaRecord);

const limitSchemaRecord: Record<string, ParamSchema> = {
  limit: {
    in: ['query'],
    errorMessage: 'limit must be a positive integer that is not greater than max: ' +
      `${config.API_LIMIT_V4}`,
    customSanitizer: {
      options: (value?: number | string): number => {
        return value !== undefined ? +value : config.API_LIMIT_V4;
      },
    },
    custom: {
      options: (value: number) => {
        // Custom validator to ensure the value is a positive integer
        if (value <= 0 || value > config.API_LIMIT_V4 || !Number.isInteger(value)) {
          throw new Error(`limit must be a positive integer that is not greater than max: ${config.API_LIMIT_V4}`);
        }
        return true;
      },
    },
  },
};

const paginationSchemaRecord: Record<string, ParamSchema> = {
  page: {
    in: ['query'],
    optional: true,
    isInt: {
      options: { gt: 0 },
    },
    errorMessage: 'page must be a non-negative integer',
  },
};

const createdBeforeOrAtSchemaRecord: Record<string, ParamSchema> = {
  createdBeforeOrAtHeight: {
    in: ['query'],
    optional: true,
    isInt: {
      options: { gt: -1 },
    },
    errorMessage: 'createdBeforeOrAtHeight must be a non-negative integer',
  },
  createdBeforeOrAt: {
    in: ['query'],
    optional: true,
    customSanitizer: {
      options: sanitizeDateToISO,
    },
    isISO8601: true,
  },
};

const effectiveBeforeOrAtSchemaRecord: Record<string, ParamSchema> = {
  effectiveBeforeOrAtHeight: {
    in: ['query'],
    optional: true,
    isInt: {
      options: { gt: -1 },
    },
    errorMessage: 'effectiveBeforeOrAtHeight must be a non-negative integer',
  },
  effectiveBeforeOrAt: {
    in: ['query'],
    optional: true,
    isISO8601: true,
  },
};

const createdOnOrAfterSchemaRecord: Record<string, ParamSchema> = {
  createdOnOrAfterHeight: {
    in: ['query'],
    optional: true,
    isInt: {
      options: { gt: -1 },
    },
    errorMessage: 'createdOnOrAfterHeight must be a non-negative integer',
  },
  createdOnOrAfter: {
    in: ['query'],
    optional: true,
    customSanitizer: {
      options: sanitizeDateToISO,
    },
    isISO8601: true,
  },
};

const transferBetweenSchemaRecord: Record<string, ParamSchema> = {
  ...createdBeforeOrAtSchemaRecord,
  sourceAddress: {
    in: ['params', 'query'],
    isString: true,
  },
  sourceSubaccountNumber: {
    in: ['params', 'query'],
    isInt: {
      options: { gt: -1, lt: MAX_PARENT_SUBACCOUNTS * CHILD_SUBACCOUNT_MULTIPLIER + 1 },
    },
    errorMessage: 'subaccountNumber must be a non-negative integer less than 128001',
  },
  recipientAddress: {
    in: ['params', 'query'],
    isString: true,
  },
  recipientSubaccountNumber: {
    in: ['params', 'query'],
    isInt: {
      options: { gt: -1, lt: MAX_PARENT_SUBACCOUNTS * CHILD_SUBACCOUNT_MULTIPLIER + 1 },
    },
    errorMessage: 'subaccountNumber must be a non-negative integer less than 128001',
  },
};

export const CheckLimitSchema = checkSchema(limitSchemaRecord);

export const CheckPaginationSchema = checkSchema(paginationSchemaRecord);

export const CheckSubaccountHistoricalFundingSchema = checkSchema({
  ...limitSchemaRecord,
  ...effectiveBeforeOrAtSchemaRecord,
  address: {
    in: ['params'],
    isString: true,
  },
  subaccountNumber: {
    in: ['params'],
    isInt: {
      options: { gt: -1, lt: MAX_PARENT_SUBACCOUNTS * CHILD_SUBACCOUNT_MULTIPLIER + 1 },
    },
    errorMessage: 'subaccountNumber must be a non-negative integer less than 128001',
  },
});

export const CheckLimitAndCreatedBeforeOrAtSchema = checkSchema({
  ...limitSchemaRecord,
  ...createdBeforeOrAtSchemaRecord,
});

export const CheckLimitAndCreatedBeforeOrAtAndOnOrAfterSchema = checkSchema({
  ...limitSchemaRecord,
  ...createdBeforeOrAtSchemaRecord,
  ...createdOnOrAfterSchemaRecord,
});

export const CheckEffectiveBeforeOrAtSchema = checkSchema({
  ...effectiveBeforeOrAtSchemaRecord,
});

const checkTickerParamSchema: ParamSchema = {
  in: 'params',
  isString: true,
  custom: {
    options: perpetualMarketRefresher.isValidPerpetualMarketTicker,
    errorMessage: 'ticker must be a valid ticker (BTC-USD, etc)',
  },
};

const checkTickerOptionalQuerySchema: ParamSchema = {
  ...checkTickerParamSchema,
  in: 'query',
  optional: true,
};

export const CheckTickerParamSchema = checkSchema({
  ticker: checkTickerParamSchema,
});

export const CheckTickerOptionalQuerySchema = checkSchema({
  ticker: checkTickerOptionalQuerySchema,
});

export const CheckHistoricalBlockTradingRewardsSchema = checkSchema({
  ...checkAddressSchemaRecord,
  ...limitSchemaRecord,
  startingBeforeOrAt: {
    in: ['query'],
    optional: true,
    isISO8601: true,
  },
  startingBeforeOrAtHeight: {
    in: ['query'],
    optional: true,
    isInt: {
      options: { gt: -1 },
    },
    errorMessage: 'startingBeforeOrAtHeight must be a non-negative integer',
  },
});

export const CheckTransferBetweenSchema = checkSchema(transferBetweenSchemaRecord);
