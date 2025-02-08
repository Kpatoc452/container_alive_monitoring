CREATE TABLE containers(
    id SERIAL PRIMARY KEY,
    address VARCHAR(25) NOT NULL,
    last_success_ping TIMESTAMP,
    last_ping TIMESTAMP,
);