CREATE TABLE IF NOT EXISTS addresses (
    id SERIAL PRIMARY KEY,
    ip INET NOT NULL UNIQUE,
    response_time text,
    last_successful_ping TIMESTAMPTZ NOT NULL
);

INSERT INTO addresses (ip, response_time, last_successful_ping) VALUES ('0.0.0.0', 10, '2001-10-05');
INSERT INTO addresses (ip, response_time, last_successful_ping) VALUES ('0.0.0.1', 11, '2001-10-05');
INSERT INTO addresses (ip, response_time, last_successful_ping) VALUES ('0.0.0.2', 12, '2001-10-05');
INSERT INTO addresses (ip, response_time, last_successful_ping) VALUES ('0.0.0.3', 13, '2001-10-05');