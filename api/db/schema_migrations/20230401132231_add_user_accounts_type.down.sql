BEGIN;
  ALTER TABLE IF EXISTS user_accounts DROP COLUMN IF EXISTS type;
  DROP TYPE IF EXISTS user_account_type;
COMMIT;