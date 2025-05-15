BEGIN;

CREATE TABLE pending_confirmation_emails (
    id SERIAL PRIMARY KEY,
    to_address VARCHAR(320) NOT NULL,
    token UUID NOT NULL,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    attempts INT NOT NULL DEFAULT 0,
    next_try_after TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

COMMIT;