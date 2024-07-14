CREATE TABLE IF NOT EXISTS balances
(
    id        SERIAL PRIMARY KEY,
    user_id   INT REFERENCES users (id) NOT NULL UNIQUE,
    current   FLOAT                     NOT NULL DEFAULT 0,
    withdrawn FLOAT                     NOT NULL DEFAULT 0
);