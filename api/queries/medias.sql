-- name: FindRawMediaCheminByID :one
SELECT chemin_media
FROM t_medias
WHERE id_media = $1
LIMIT 1;


-- name: UpdateMediaTitle :exec
UPDATE t_medias
SET titre_media = sqlc.arg(title),
    chemin_media = jsonb_set(chemin_media::jsonb, '{0,title}', to_jsonb(sqlc.arg(json_title)::text))::text,
    date_maj = NOW()
WHERE id_media = sqlc.arg(id);

-- name: DuplicateMedia :one
INSERT INTO t_medias(titre_media, chemin_media, date_creation)
SELECT m.titre_media, m.chemin_media, m.date_creation
FROM t_medias m
WHERE m.id_media = sqlc.arg(source_id)
RETURNING id_media;

-- name: GetMediaIdsByMonuLieu :many
SELECT media_monu_lieu_id FROM cor_medias_monu_lieu WHERE monument_lieu_id = sqlc.arg(id);

-- name: GetMediaIdsByMobImg :many
SELECT media_mob_img_id FROM cor_medias_mob_img WHERE mobilier_image_id = sqlc.arg(id);

-- name: GetMediaIdsByPersMo :many
SELECT media_pers_mo_id FROM cor_medias_pers_mo WHERE pers_morale_id = sqlc.arg(id);

-- name: GetMediaIdsByPersPhy :many
SELECT media_pers_phy_id FROM cor_medias_pers_phy WHERE pers_physique_id = sqlc.arg(id);

-- name: CreateNewMedia :one
INSERT INTO t_medias(titre_media, chemin_media, date_creation)
VALUES (sqlc.arg(title), sqlc.arg(chemin_media), sqlc.arg(date_creation))
RETURNING id_media;