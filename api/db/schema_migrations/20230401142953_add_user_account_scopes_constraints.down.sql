BEGIN;
  ALTER TABLE IF EXISTS user_accounts ALTER COLUMN scopes DROP NOT NULL;
  ALTER TABLE IF EXISTS user_accounts ALTER COLUMN scopes DROP DEFAULT;
COMMIT;