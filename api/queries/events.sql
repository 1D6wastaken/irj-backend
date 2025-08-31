-- name: ContributorRegistrationEvent :exec
INSERT INTO t_app_events (type, user_id)
VALUES ('contributor_registration', sqlc.arg(id));

-- name: ContributorValidationEvent :exec
INSERT INTO t_app_events (type, user_id, admin_id)
VALUES ('contributor_validation', sqlc.arg(user_id), sqlc.arg(admin_id));

-- name: ContributorRejectionEvent :exec
INSERT INTO t_app_events (type, user_id, admin_id)
VALUES ('contributor_rejection', sqlc.arg(user_id), sqlc.arg(admin_id));

-- name: AccountDeletionEvent :exec
INSERT INTO t_app_events (type, user_id)
VALUES ('account_deletion', sqlc.arg(id));

-- name: DocumentSubmissionEvent :exec
INSERT INTO t_app_events (type, user_id, document_id, comment)
VALUES ('document_submission', sqlc.arg(user_id), sqlc.arg(document_id), sqlc.arg(comment));

-- name: DocumentValidationEvent :exec
INSERT INTO t_app_events (type, document_id, admin_id, comment)
VALUES ('document_validation', sqlc.arg(document_id), sqlc.arg(admin_id), sqlc.arg(comment));

-- name: DocumentRejectionEvent :exec
INSERT INTO t_app_events (type, document_id, admin_id, comment)
VALUES ('document_rejection', sqlc.arg(document_id), sqlc.arg(admin_id), sqlc.arg(comment));

-- name: DocumentUpdateEvent :exec
INSERT INTO t_app_events (type, user_id, document_id, comment)
VALUES ('document_update', sqlc.arg(user_id), sqlc.arg(document_id), sqlc.arg(comment));
