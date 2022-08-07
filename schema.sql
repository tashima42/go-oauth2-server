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

INSERT INTO
  clients (name, client_id, client_secret, redirect_uri)
VALUES
  (
    'client name',
    'client1',
    '$2b$10$P9PjYWou7PU.pDA3sx3DwuW1ny902LV13LVZsZGHlahuOUbsOPuBO',
    'https://sp-dev.tbxnet.com/v2/auth/oauth2/assert'
  );

INSERT INTO
  user_accounts (username, password, country, subscriber_id)
VALUES
  (
    'user1@example.com',
    '$2b$10$P9PjYWou7PU.pDA3sx3DwuW1ny902LV13LVZsZGHlahuOUbsOPuBO',
    'AR',
    'subscriber1'
  );
