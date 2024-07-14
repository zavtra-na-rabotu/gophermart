CREATE TABLE IF NOT EXISTS users
(
    id       SERIAL PRIMARY KEY,
    login    VARCHAR(255) UNIQUE,
    password VARCHAR(72)
);