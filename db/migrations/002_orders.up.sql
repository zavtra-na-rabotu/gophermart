CREATE TABLE IF NOT EXISTS orders
(
    id          SERIAL PRIMARY KEY        NOT NULL,
    number      VARCHAR(255) UNIQUE       NOT NULL,
    status      VARCHAR(50)               NOT NULL,
    accrual     FLOAT                     NOT NULL DEFAULT 0,
    user_id     INT REFERENCES users (id) NOT NULL,
    uploaded_at TIMESTAMP WITH TIME ZONE  NOT NULL DEFAULT CURRENT_TIMESTAMP
);