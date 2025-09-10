-- name: GetDepartements :many
SELECT d.id_departement AS id,
       d.nom_departement AS name,
       d.id_region AS region_id,
       r.nom_region AS region_name,
       r.id_pays AS pays_id,
       p.nom_pays AS pays_name
FROM loc_departements d
         JOIN loc_regions r ON d.id_region = r.id_region
         JOIN loc_pays p ON r.id_pays = p.id_pays
WHERE d.nom_departement != '' AND d.nom_departement IS NOT NULL;

-- name: GetRegions :many
SELECT r.id_region AS id,
       r.nom_region AS name,
       r.id_pays AS pays_id,
       p.nom_pays AS pays_name
FROM loc_regions r
         JOIN loc_pays p ON r.id_pays = p.id_pays
WHERE r.nom_region != '' AND r.nom_region IS NOT NULL;

-- name: GetPays :many
SELECT id_pays AS id, nom_pays AS name FROM loc_pays WHERE nom_pays != '' AND nom_pays IS NOT NULL;

-- name: GetSiecles :many
SELECT id_siecle AS id, siecle_list AS name FROM bib_siecle WHERE siecle_list != '' AND siecle_list IS NOT NULL ORDER BY id_siecle ASC;

-- name: GetEtatsConservation :many
SELECT id_etat_conservation AS id, etat_conservation_type AS name FROM bib_etats_conservation WHERE etat_conservation_type != '' AND etat_conservation_type IS NOT NULL;

-- name: GetNaturesMonu :many
SELECT id_monu_lieu_nature AS id, monu_lieu_nature_type AS name FROM bib_monu_lieu_natures WHERE monu_lieu_nature_type != '' AND monu_lieu_nature_type IS NOT NULL;

-- name: GetMateriaux :many
SELECT id_materiau AS id, materiau_type AS name FROM bib_materiaux WHERE materiau_type != '' AND materiau_type IS NOT NULL;

-- name: GetNaturesMob :many
SELECT id_nature AS id, nature_type AS name FROM bib_mob_img_natures WHERE nature_type != '' AND nature_type IS NOT NULL;

-- name: GetTechniquesMob :many
SELECT id_technique AS id, technique_type AS name FROM bib_mob_img_techniques WHERE technique_type != '' AND technique_type IS NOT NULL;

-- name: GetNaturesPersonnesMorales :many
SELECT id_pers_mo_nature AS id, pers_mo_nature_type AS name FROM bib_pers_mo_natures WHERE pers_mo_nature_type != '' AND pers_mo_nature_type IS NOT NULL;

-- name: GetProfessions :many
SELECT id_profession AS id, profession_type AS name FROM bib_pers_phy_professions WHERE profession_type != '' AND profession_type IS NOT NULL;

-- name: GetDeplacements :many
SELECT id_mode_deplacement AS id, mode_deplacement_type AS name FROM bib_pers_phy_modes_deplacements WHERE mode_deplacement_type != '' AND mode_deplacement_type IS NOT NULL;

-- name: GetHistoricalPeriods :many
SELECT id_periode_historique AS id, periode_historique_type AS name FROM bib_pers_phy_periodes_historiques WHERE periode_historique_type != '' AND periode_historique_type IS NOT NULL;

-- name: GetThemes :many
SELECT id_theme AS id, theme_type AS name FROM t_themes WHERE theme_type != '' AND theme_type IS NOT NULL;

-- name: SearchCommunesPaginated :many
SELECT
    c.id_commune AS id,
    c.nom_commune AS commune_name,
    c.id_departement AS departement_id,
    d.nom_departement AS departement_name,
    d.id_region AS region_id,
    r.nom_region AS region_name,
    r.id_pays AS pays_id,
    p.nom_pays AS pays_name
FROM loc_communes c
         JOIN loc_departements d ON c.id_departement = d.id_departement
         JOIN loc_regions r ON d.id_region = r.id_region
         JOIN loc_pays p ON r.id_pays = p.id_pays
WHERE c.nom_commune != '' AND c.nom_commune IS NOT NULL AND c.nom_commune ILIKE $1
ORDER BY c.nom_commune ASC
LIMIT $2 OFFSET $3;