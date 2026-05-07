import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex.schema.alterTable('user_complaints', (table) => {
    table.string('walletAddress').notNullable().defaultTo('');
  });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.alterTable('user_complaints', (table) => {
    table.dropColumn('walletAddress');
  });
}
