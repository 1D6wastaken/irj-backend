-- name: CreateResetToken :exec
INSERT INTO t_password_resets (user_id, token)
VALUES (sqlc.arg(id), sqlc.arg(token));

-- name: GetResetPasswordByToken :one
SELECT user_id, token, created_at
FROM t_password_resets
WHERE token = $1;

-- name: DeletePasswordResetByID :exec
DELETE
FROM t_password_resets
WHERE user_id = $1;

-- name: DeletePasswordResetByToken :exec
DELETE
FROM t_password_resets
WHERE token = $1;

-- name: DeleteExpiredPasswordReset :exec
DELETE
FROM t_password_resets
WHERE created_at < NOW() - INTERVAL '30 minutes';
