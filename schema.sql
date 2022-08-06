CREATE TABLE IF NOT EXISTS clients(
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  client_id TEXT NOT NULL UNIQUE,
  client_secret TEXT NOT NULL,
  redirect_uri TEXT NOT NULL,
  grants TEXT
);

CREATE TABLE IF NOT EXISTS user_accounts(
  id SERIAL PRIMARY KEY,
  username TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL UNIQUE,
  country TEXT NOT NULL,
  subscriber_id TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS tokens(
  id SERIAL PRIMARY KEY,
  access_token TEXT NOT NULL UNIQUE,
  access_token_expires_at TEXT NOT NULL,
  refresh_token_expires_at TEXT NOT NULL,
  refresh_token TEXT NOT NULL UNIQUE,
  client_id SERIAL NOT NULL,
  user_account_id SERIAL NOT NULL,
  active BOOLEAN NOT NULL DEFAULT true,
  FOREIGN KEY (client_id) REFERENCES clients (id),
  FOREIGN KEY (user_account_id) REFERENCES user_accounts (id)
);

CREATE TABLE IF NOT EXISTS authorization_codes(
  id SERIAL PRIMARY KEY,
  code TEXT NOT NULL UNIQUE,
  expires_at TEXT NOT NULL,
  redirect_uri TEXT NOT NULL,
  client_id SERIAL NOT NULL,
  user_account_id SERIAL NOT NULL,
  active BOOLEAN NOT NULL DEFAULT true,
  FOREIGN KEY (client_id) REFERENCES clients (id),
  FOREIGN KEY (user_account_id) REFERENCES user_accounts (id)
);
