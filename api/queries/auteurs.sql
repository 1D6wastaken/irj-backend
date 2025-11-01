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

-- name: GetAuteurIDByPersPhyID :one
SELECT auteur_fiche_pers_phy_id
FROM cor_auteur_fiche_pers_phy
WHERE pers_physique_id = $1;

-- name: GetAuteurIDByPersMoID :one
SELECT auteur_fiche_pers_mo_id
FROM cor_auteur_fiche_pers_mo
WHERE pers_morale_id = $1;

-- name: GetAuteurIDByMonuLieu :one
SELECT auteur_fiche_monu_lieu_id
FROM cor_auteur_fiche_monu_lieu
WHERE monument_lieu_id = $1;

-- name: GetAuteurIDByMobImg :one
SELECT auteur_fiche_mob_img_id
FROM cor_auteur_fiche_mob_img
WHERE mobilier_image_id = $1;