BEGIN;

ALTER TABLE user_subscriptions ADD COLUMN confirmed BOOLEAN DEFAULT FALSE;

COMMIT;