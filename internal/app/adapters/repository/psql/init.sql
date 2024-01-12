-- users

CREATE TABLE IF NOT EXISTS users (
	_id SERIAL PRIMARY KEY,
	uuid VARCHAR(36) NOT NULL,
	login VARCHAR(250) NOT NULL,
	password VARCHAR(250) NOT NULL,
	is_active BOOLEAN DEFAULT true
);

CREATE UNIQUE INDEX IF NOT EXISTS login_idx ON users (login);


-- orders

CREATE TABLE IF NOT EXISTS orders (
	_id SERIAL PRIMARY KEY,
	uuid VARCHAR(250) NOT NULL,
	status VARCHAR(25) NOT NULL,
	accrual NUMERIC DEFAULT 0,
	uploaded_at TIMESTAMPTZ NOT NULL,
	user_id VARCHAR(36) NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS uuid_idx ON orders (uuid);


-- balance

CREATE TABLE IF NOT EXISTS balance (
	_id SERIAL PRIMARY KEY,
	user_id VARCHAR(36) NOT NULL,
	current NUMERIC DEFAULT 0 CHECK (current >= 0),
	withdrawn NUMERIC DEFAULT 0  CHECK (withdrawn >= 0),
	updated_at TIMESTAMPTZ NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS user_id_idx ON balance (user_id);


-- withdrawals

CREATE TABLE IF NOT EXISTS withdrawals (
	_id SERIAL PRIMARY KEY,
	order_id VARCHAR(250) NOT NULL,
	user_id VARCHAR(36) NOT NULL,
	sum NUMERIC DEFAULT 0,
	processed_at TIMESTAMPTZ NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS user_id_idx ON balance (user_id);
