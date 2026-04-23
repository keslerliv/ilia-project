-- USER SCHEMA
CREATE TABLE
    IF NOT EXISTS "address" (
        id SERIAL PRIMARY KEY,
        user_id INT NOT NULL,
        FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE
    );