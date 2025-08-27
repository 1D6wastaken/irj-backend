-- name: FindRawMediaCheminByID :one
SELECT chemin_media
FROM t_medias
WHERE id_media = $1
LIMIT 1;


-- name: CreateNewMedia :one
INSERT INTO t_medias(titre_media, chemin_media, date_creation)
VALUES (sqlc.arg(title), sqlc.arg(chemin_media), sqlc.arg(date_creation))
RETURNING id_media;