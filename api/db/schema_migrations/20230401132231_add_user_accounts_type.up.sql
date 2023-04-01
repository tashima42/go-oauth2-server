BEGIN;
  CREATE TYPE user_account_type AS ENUM ('dev', 'user');
  ALTER TABLE IF EXISTS user_accounts ADD COLUMN IF NOT EXISTS type user_account_type NOT NULL DEFAULT 'user';
COMMIT;