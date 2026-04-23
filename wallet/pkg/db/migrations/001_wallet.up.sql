CREATE TABLE
    IF NOT EXISTS "wallets" (
        id SERIAL PRIMARY KEY,
        user_id INT UNIQUE NOT NULL,
        balance INT NOT NULL DEFAULT 0,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

CREATE TABLE
    IF NOT EXISTS "wallet_transactions" (
        id SERIAL PRIMARY KEY,
        user_id INT NOT NULL,
        amount INT NOT NULL DEFAULT 0,
        action VARCHAR(255) NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );