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

-- name: DocumentSubmissionValidationEvent :exec
INSERT INTO t_app_events (type, document_id, admin_id, comment)
VALUES ('document_submission_validation', sqlc.arg(document_id), sqlc.arg(admin_id), sqlc.arg(comment));

-- name: DocumentSubmissionRejectionEvent :exec
INSERT INTO t_app_events (type, document_id, admin_id, comment)
VALUES ('document_submission_rejection', sqlc.arg(document_id), sqlc.arg(admin_id), sqlc.arg(comment));

-- name: DocumentUpdateEvent :exec
INSERT INTO t_app_events (type, user_id, document_id, comment)
VALUES ('document_update', sqlc.arg(user_id), sqlc.arg(document_id), sqlc.arg(comment));

-- name: DocumentUpdateValidationEvent :exec
INSERT INTO t_app_events (type, document_id, admin_id, comment)
VALUES ('document_update_validation', sqlc.arg(document_id), sqlc.arg(admin_id), sqlc.arg(comment));

-- name: DocumentUpdateRejectionEvent :exec
INSERT INTO t_app_events (type, document_id, admin_id, comment)
VALUES ('document_update_rejection', sqlc.arg(document_id), sqlc.arg(admin_id), sqlc.arg(comment));

-- name: UpdateDocumentIDAfterUpdate :exec
UPDATE t_app_events
SET document_id = sqlc.arg(new_doc_id)
WHERE document_id = sqlc.arg(parent_doc_id);

-- name: GetHistoryByUserID :many
WITH user_docs AS (
    SELECT DISTINCT document_id
    FROM t_app_events
    WHERE t_app_events.user_id = sqlc.arg(user_id)
      AND type IN ('document_submission', 'document_update', 'document_submission_validation', 'document_submission_rejection', 'document_update_validation', 'document_update_rejection')
),
     user_related_events AS (
         SELECT e.*
         FROM t_app_events e
         WHERE e.user_id = sqlc.arg(user_id)
           AND type IN ('document_submission', 'document_update', 'document_submission_validation', 'document_submission_rejection', 'document_update_validation', 'document_update_rejection')

         UNION

         SELECT e.*
         FROM t_app_events e
                  JOIN user_docs ud ON ud.document_id = e.document_id
         WHERE e.admin_id IS NOT NULL
           AND type IN ('document_submission', 'document_update', 'document_submission_validation', 'document_submission_rejection', 'document_update_validation', 'document_update_rejection')

     ),
     scored_events AS (
         SELECT
             e.*,
             CASE
                 WHEN e.type IN (
                                 'document_submission_validation',
                                 'document_submission_rejection',
                                 'document_update_validation',
                                 'document_update_rejection'
                     ) THEN 2
                 ELSE 1
                 END AS priority
         FROM user_related_events e
     ),
     ranked AS (
         SELECT
             *,
             ROW_NUMBER() OVER (
                 PARTITION BY document_id
                 ORDER BY priority DESC, date DESC
                 ) AS rn
         FROM scored_events
     ),
     selected AS (
         SELECT *
         FROM ranked
         WHERE rn = 1
     ),
     paginated AS (
         SELECT
             *,
             COUNT(*) OVER() AS total_count
         FROM selected
         ORDER BY date DESC
     )
SELECT *
FROM paginated
LIMIT sqlc.arg(limit_param) OFFSET sqlc.arg(offset_param);


-- name: GetAllHistory :many
WITH user_docs AS (
    SELECT DISTINCT document_id
    FROM t_app_events
    WHERE type IN (
                   'document_submission', 'document_update',
                   'document_submission_validation', 'document_submission_rejection',
                   'document_update_validation', 'document_update_rejection'
        )
),

-- ✅ Récupère l'événement utilisateur qui a créé ou mis à jour la fiche
     creator_events AS (
         SELECT DISTINCT ON (document_id)
             document_id,
             user_id AS creator_user_id
         FROM t_app_events
         WHERE type IN ('document_submission','document_update')
         ORDER BY document_id, date ASC   -- le 1er event utilisateur
     ),

     user_related_events AS (
         SELECT e.*
         FROM t_app_events e
         WHERE type IN (
                        'document_submission', 'document_update',
                        'document_submission_validation', 'document_submission_rejection',
                        'document_update_validation', 'document_update_rejection'
             )

         UNION

         SELECT e.*
         FROM t_app_events e
                  JOIN user_docs ud ON ud.document_id = e.document_id
         WHERE e.admin_id IS NOT NULL
           AND type IN (
                        'document_submission', 'document_update',
                        'document_submission_validation', 'document_submission_rejection',
                        'document_update_validation', 'document_update_rejection'
             )
     ),

     scored_events AS (
         SELECT
             e.*,
             CASE
                 WHEN e.type IN (
                                 'document_submission_validation',
                                 'document_submission_rejection',
                                 'document_update_validation',
                                 'document_update_rejection'
                     ) THEN 2
                 ELSE 1
                 END AS priority
         FROM user_related_events e
     ),

     ranked AS (
         SELECT
             *,
             ROW_NUMBER() OVER (
                 PARTITION BY document_id
                 ORDER BY priority DESC, date DESC
                 ) AS rn
         FROM scored_events
     ),

     selected AS (
         SELECT *
         FROM ranked
         WHERE rn = 1
     ),

     paginated AS (
         SELECT
             s.*,
             ce.creator_user_id,   -- ✅ récupère l'user d'origine
             u.prenom  AS user_prenom,
             u.nom     AS user_nom,
             u.email   AS user_email,
             a.prenom  AS admin_prenom,
             a.nom     AS admin_nom,
             a.email   AS admin_email,
             COUNT(*) OVER() AS total_count
         FROM selected s
                  LEFT JOIN creator_events ce ON ce.document_id = s.document_id
                  LEFT JOIN t_app_users u ON u.id = ce.creator_user_id
                  LEFT JOIN t_app_users a ON a.id = s.admin_id
         ORDER BY s.date DESC
     )

SELECT *
FROM paginated
LIMIT sqlc.arg(limit_param) OFFSET sqlc.arg(offset_param);