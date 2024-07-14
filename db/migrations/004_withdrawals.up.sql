CREATE TABLE IF NOT EXISTS withdrawals
(
    id           SERIAL PRIMARY KEY,
    user_id      INT REFERENCES users (id) NOT NULL,
    order_number VARCHAR(255)              NOT NULL,
    sum          FLOAT                     NOT NULL,
    processed_at TIMESTAMP WITH TIME ZONE  NOT NULL DEFAULT CURRENT_TIMESTAMP
);