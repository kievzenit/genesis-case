BEGIN;

ALTER TABLE user_subscriptions DROP COLUMN confirmed;

COMMIT;