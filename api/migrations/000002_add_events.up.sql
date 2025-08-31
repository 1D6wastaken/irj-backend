BEGIN;


CREATE TYPE event_type AS ENUM  (
    'contributor_registration',
    'contributor_validation',
    'contributor_rejection',
    'account_deletion',
    'document_submission',
    'document_validation',
    'document_rejection',
    'document_update'
    );

CREATE TABLE IF NOT EXISTS t_app_events
(
    id SERIAL PRIMARY KEY,
    type event_type NOT NULL,
    user_id CHAR(36) NOT NULL,
    admin_id CHAR(36),
    document_id INTEGER,
    comment TEXT,
    date TIMESTAMPTZ       NOT NULL DEFAULT NOW()
);

COMMIT;