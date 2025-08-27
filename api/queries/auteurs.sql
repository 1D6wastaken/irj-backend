-- name: CreateAuteur :one
INSERT INTO bib_auteurs (auteur_fiche_nom)
VALUES (sqlc.arg(name))
RETURNING id_auteur_fiche;

-- name: GetAuteurByName :one
SELECT *
FROM bib_auteurs
WHERE auteur_fiche_nom = $1;