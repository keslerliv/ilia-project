CREATE TABLE
    IF NOT EXISTS "wallets" (
        id SERIAL PRIMARY KEY,
        user_id INT UNIQUE NOT NULL,
        balance INT NOT NULL DEFAULT 0,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        FOREIGN KEY (user_id) REFERENCES "users" (id) ON DELETE CASCADE
    );