BEGIN;

CREATE TABLE frequencies (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE user_subscriptions (
    id SERIAL PRIMARY KEY,
    token UUID NOT NULL UNIQUE,
    email VARCHAR(320) NOT NULL,
    city VARCHAR(100) NOT NULL,
    frequency_id INT NOT NULL REFERENCES frequencies(id)
);

CREATE UNIQUE INDEX idx_user_subscriptions_token ON user_subscriptions(token);

INSERT INTO frequencies (name) VALUES
    ('hourly'),
    ('daily');

COMMIT;