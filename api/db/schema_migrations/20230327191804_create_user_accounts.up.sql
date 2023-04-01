CREATE TABLE IF NOT EXISTS user_accounts(
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  username TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL UNIQUE
);
