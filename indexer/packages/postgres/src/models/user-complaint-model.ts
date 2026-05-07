import { Model } from 'objection';

import UpsertQueryBuilder from '../query-builders/upsert';

export default class UserComplaintModel extends Model {
  static get tableName() {
    return 'user_complaints';
  }

  static get idColumn() {
    return 'id';
  }

  static relationMappings = {};

  static get jsonSchema() {
    return {
      type: 'object',
      required: ['id', 'message'],
      properties: {
        id: { type: 'string', format: 'uuid' },
        message: { type: 'string' },
        email: { type: ['string', 'null'] },
        createdAt: { type: 'string', format: 'date-time' },
      },
    };
  }

  QueryBuilderType!: UpsertQueryBuilder<this>;

  id!: string;

  message!: string;

  email?: string;

  createdAt!: string;
}
