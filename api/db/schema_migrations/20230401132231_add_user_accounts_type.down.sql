BEGIN;
  DROP TYPE IF EXISTS user_account_type;
  ALTER TABLE IF EXISTS user_accounts DROP COLUMN IF EXISTS type;
COMMIT;