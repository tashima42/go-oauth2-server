CREATE TABLE IF NOT EXISTS authorization_codes(
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  code TEXT NOT NULL UNIQUE,
  expires_at TEXT NOT NULL,
  redirect_uri TEXT NOT NULL,
  client_id uuid NOT NULL,
  user_account_id uuid NOT NULL,
  active BOOLEAN NOT NULL DEFAULT true,
  FOREIGN KEY (client_id) REFERENCES clients (id),
  FOREIGN KEY (user_account_id) REFERENCES user_accounts (id)
);