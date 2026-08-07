import * as Knex from 'knex';

// Note: `20260507000000_create_user_complaints_table` was edited after it had already been
// applied, so it now creates `walletAddress` itself. On the live database this migration is
// still the one that added the column; on any database migrated from scratch the column
// already exists by the time this runs, and adding it unconditionally fails with
// `column "walletAddress" of relation "user_complaints" already exists`.
//
// Guarding on `hasColumn` makes both paths work without rewriting migration history again.
export async function up(knex: Knex): Promise<void> {
  const hasColumn: boolean = await knex.schema.hasColumn('user_complaints', 'walletAddress');
  if (hasColumn) {
    return;
  }
  return knex.schema.alterTable('user_complaints', (table) => {
    table.string('walletAddress').notNullable().defaultTo('');
  });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.alterTable('user_complaints', (table) => {
    table.dropColumn('walletAddress');
  });
}
