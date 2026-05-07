import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex.schema.createTable('user_complaints', (table) => {
    table.uuid('id').primary();
    table.text('message').notNullable();
    table.string('email').nullable();
    table.timestamp('createdAt').defaultTo(knex.fn.now()).notNullable();
  });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTable('user_complaints');
}
