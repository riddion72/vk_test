CREATE TABLE IF NOT EXISTS addresses (
    id SERIAL PRIMARY KEY,
    ip INET NOT NULL UNIQUE,
    response_time NUMERIC(10,3)
    last_successful_ping TIMESTAMPTZ NOT NULL,
);