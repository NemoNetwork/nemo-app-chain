import { QueryBuilder } from 'objection';
import { v4 as uuidv4 } from 'uuid';

import { DEFAULT_POSTGRES_OPTIONS } from '../constants';
import {
  setupBaseQuery,
  verifyAllRequiredFields,
} from '../helpers/stores-helpers';
import Transaction from '../helpers/transaction';
import UserComplaintModel from '../models/user-complaint-model';
import {
  Options,
  Ordering,
  QueryableField,
  QueryConfig,
  UserComplaintColumns,
  UserComplaintCreateObject,
  UserComplaintFromDatabase,
  UserComplaintQueryConfig,
} from '../types';

export async function findAll(
  {
    createdBeforeOrAt,
    limit,
  }: UserComplaintQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<UserComplaintFromDatabase[]> {
  verifyAllRequiredFields(
    {
      createdBeforeOrAt,
      limit,
    } as QueryConfig,
    requiredFields,
  );

  let baseQuery: QueryBuilder<UserComplaintModel> = setupBaseQuery<UserComplaintModel>(
    UserComplaintModel,
    options,
  );

  if (createdBeforeOrAt !== undefined) {
    baseQuery = baseQuery.where(UserComplaintColumns.createdAt, '<=', createdBeforeOrAt);
  }

  if (options.orderBy !== undefined) {
    for (const [column, order] of options.orderBy) {
      baseQuery = baseQuery.orderBy(column, order);
    }
  } else {
    baseQuery = baseQuery.orderBy(UserComplaintColumns.createdAt, Ordering.DESC);
  }

  if (limit !== undefined) {
    baseQuery = baseQuery.limit(limit);
  }

  return baseQuery.returning('*');
}

export async function create(
  complaintToCreate: Omit<UserComplaintCreateObject, 'id'>,
  options: Options = { txId: undefined },
): Promise<UserComplaintFromDatabase> {
  return UserComplaintModel.query(
    Transaction.get(options.txId),
  ).insert({
    id: uuidv4(),
    ...complaintToCreate,
  }).returning('*');
}
