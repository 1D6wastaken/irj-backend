-- name: CreateAuteur :one
INSERT INTO bib_auteurs (auteur_fiche_nom, user_id)
VALUES (sqlc.arg(name), sqlc.arg(user_id))
RETURNING id_auteur_fiche;

-- name: GetAuteurByName :one
SELECT *
FROM bib_auteurs
WHERE auteur_fiche_nom = $1;

-- name: GetAuteurByUserID :one
SELECT *
FROM bib_auteurs
WHERE user_id = $1;