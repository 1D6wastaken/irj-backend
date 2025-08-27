-- name: CreateUser :exec
INSERT INTO t_app_users (id, prenom, nom, email, mot_de_passe, telephone, organisation, domaine, motivation)
VALUES (sqlc.arg(id), sqlc.arg(firstname), sqlc.arg(name), sqlc.arg(email), sqlc.arg(password), sqlc.arg(phone),
        sqlc.arg(organization), sqlc.arg(domain), sqlc.arg(motivation));

-- name: GetUserByID :one
SELECT prenom,
       nom,
       mot_de_passe,
       email,
       email_confirm,
       telephone,
       organisation,
       domaine,
       motivation,
       grade,
       date_creation
FROM t_app_users
WHERE id = $1;

-- name: GetUsers :many
SELECT id,
       prenom,
       nom,
       email,
       last_login
FROM t_app_users;

-- name: GetUserByEmail :one
SELECT id,
       prenom,
       nom,
       mot_de_passe,
       email,
       email_confirm,
       telephone,
       organisation,
       domaine,
       motivation,
       grade,
       date_creation
FROM t_app_users
WHERE email = $1;

-- name: UpdateLastLogin :exec
UPDATE t_app_users
SET last_login = NOW()
WHERE id = $1;

-- name: GetUsersByGrade :many
SELECT id,
       email,
       prenom,
       nom,
       organisation,
       domaine,
       date_creation,
       motivation,
       telephone
FROM t_app_users
WHERE grade = $1;

-- name: UpdatePasswordByID :exec
UPDATE t_app_users
SET mot_de_passe = sqlc.arg(password)
WHERE id = $1;

-- name: UpdateUserByID :exec
UPDATE t_app_users
SET prenom        = sqlc.arg(firstname),
    nom           = sqlc.arg(name),
    email         = sqlc.arg(email),
    email_confirm = sqlc.arg(email_confirm),
    mot_de_passe  = sqlc.arg(password),
    telephone     = sqlc.arg(phone),
    organisation  = sqlc.arg(organization),
    domaine       = sqlc.arg(domain),
    date_maj      = now()
WHERE id = $1;

-- name: DeleteUserByID :exec
DELETE
FROM t_app_users
WHERE id = $1;

-- name: ApproveUserByID :exec
UPDATE t_app_users
SET grade = 'ACTIVE'
WHERE id = $1;

-- name: ConfirmEmailUserByID :exec
UPDATE t_app_users
SET email_confirm = true
WHERE id = $1;