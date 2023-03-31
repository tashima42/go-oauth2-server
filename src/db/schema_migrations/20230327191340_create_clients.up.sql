CREATE TABLE IF NOT EXISTS clients(
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  name TEXT NOT NULL,
  client_id TEXT NOT NULL UNIQUE,
  client_secret TEXT NOT NULL UNIQUE,
  redirect_uri TEXT NOT NULL
);