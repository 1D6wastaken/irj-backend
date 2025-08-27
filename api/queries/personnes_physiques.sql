-- name: GetFilteredPersonnesPhysiques :many
SELECT p.id_pers_physique                                                                             AS id,
       p.prenom_nom_pers_phy                                                                          AS firstname,
       COALESCE(array_agg(DISTINCT bs.siecle_list) FILTER (WHERE bs.siecle_list IS NOT NULL),
                '{}')                                                                                 AS siecles,
       COALESCE(array_agg(DISTINCT bpp.profession_type) FILTER (WHERE bpp.profession_type IS NOT NULL),
                '{}')                                                                                 AS professions,
       COALESCE(array_agg(DISTINCT tm.chemin_media) FILTER (WHERE tm.chemin_media IS NOT NULL), '{}') AS medias
FROM t_pers_physiques p
         LEFT JOIN cor_siecles_pers_phy csp ON p.id_pers_physique = csp.pers_physique_id
         LEFT JOIN bib_siecle bs ON bs.id_siecle = csp.siecle_pers_phy_id

         LEFT JOIN cor_professions_pers_phy cpp ON p.id_pers_physique = cpp.pers_physique_id
         LEFT JOIN bib_pers_phy_professions bpp ON bpp.id_profession = cpp.profession_id

         LEFT JOIN cor_medias_pers_phy cmp ON p.id_pers_physique = cmp.pers_physique_id
         LEFT JOIN t_medias tm ON tm.id_media = cmp.media_pers_phy_id

         LEFT JOIN loc_communes c ON p.id_commune = c.id_commune
         LEFT JOIN loc_departements d ON c.id_departement = d.id_departement
         LEFT JOIN loc_regions r ON d.id_region = r.id_region
         LEFT JOIN loc_pays pa ON r.id_pays = pa.id_pays

         LEFT JOIN cor_auteur_fiche_pers_phy caf ON p.id_pers_physique = caf.pers_physique_id
         LEFT JOIN cor_mob_img_pers_phy cpm ON p.id_pers_physique = cpm.pers_physique_id
         LEFT JOIN cor_pers_phy_pers_mo cpmo ON p.id_pers_physique = cpmo.pers_morale_id
         LEFT JOIN cor_monu_lieu_pers_phy cmi ON p.id_pers_physique = cmi.pers_phy_id
WHERE (cardinality(sqlc.arg(siecles)::int[]) = 0 OR csp.siecle_pers_phy_id = ANY (sqlc.arg(siecles)::int[]))
   OR (cardinality(sqlc.arg(professions)::int[]) = 0 OR cpp.profession_id = ANY (sqlc.arg(professions)::int[]))
   OR (cardinality(sqlc.arg(auteurs_fiche)::int[]) = 0 OR
       caf.auteur_fiche_pers_phy_id = ANY (sqlc.arg(auteurs_fiche)::int[]))
   OR (cardinality(sqlc.arg(pers_mo)::int[]) = 0 OR cpp.pers_physique_id = ANY (sqlc.arg(pers_mo)::int[]))
   OR (cardinality(sqlc.arg(places)::int[]) = 0 OR cmi.monu_lieu_id = ANY (sqlc.arg(places)::int[]))
   OR (cardinality(sqlc.arg(furniture)::int[]) = 0 OR cpm.mobilier_image_id = ANY (sqlc.arg(furniture)::int[]))
   OR (cardinality(sqlc.arg(cities)::int[]) = 0 OR c.id_commune = ANY (sqlc.arg(cities)::int[]))
   OR (cardinality(sqlc.arg(departments)::int[]) = 0 OR d.id_departement = ANY (sqlc.arg(departments)::int[]))
   OR (cardinality(sqlc.arg(regions)::int[]) = 0 OR r.id_region = ANY (sqlc.arg(regions)::int[]))
   OR (cardinality(sqlc.arg(pays)::int[]) = 0 OR pa.id_pays = ANY (sqlc.arg(pays)::int[]))
GROUP BY p.id_pers_physique
ORDER BY p.id_pers_physique
LIMIT sqlc.arg(limit_param) OFFSET sqlc.arg(offset_param);

-- name: GetPersonnePhysiqueByID :one
SELECT p.id_pers_physique    AS id,
       p.prenom_nom_pers_phy AS firstname,
       p.date_naissance,
       p.date_deces,
       p.attestation,
       -- Periode historique
       COALESCE(array_agg(DISTINCT bph.periode_historique_type) FILTER (WHERE bph.periode_historique_type IS NOT NULL),
                '{}')        AS historical_period,
       p.bibliographie,
       p.elements_biographiques,
       p.elements_pelerinage,
       p.commutation_voeu,
       p.sources,
       p.date_creation,
       p.date_maj,
       p.publie,
       p.contributeurs,
       p.commentaires,
       -- Redacteurs (auteurs fiche)
       COALESCE(array_agg(DISTINCT baf.auteur_fiche_nom) FILTER (WHERE baf.auteur_fiche_nom IS NOT NULL),
                '{}')        AS redacteurs,
       -- Commune
       c.nom_commune         AS commune,
       -- Département
       d.nom_departement     AS departement,
       -- Région
       r.nom_region          AS region,
       -- Pays
       pa.nom_pays           AS pays,
       -- Travels
       COALESCE(array_agg(DISTINCT bmd.mode_deplacement_type) FILTER (WHERE bmd.mode_deplacement_type IS NOT NULL),
                '{}')        AS travels,
       -- Professions
       COALESCE(array_agg(DISTINCT bpp.profession_type) FILTER (WHERE bpp.profession_type IS NOT NULL),
                '{}')        AS professions,
       p.nature_evenement,
       -- Médias
       COALESCE(
                       jsonb_agg(
                       DISTINCT jsonb_build_object(
                               'id', tm.id_media,
                               'titre', tm.titre_media
                                )
                                ) FILTER (
                           WHERE tm.chemin_media IS NOT NULL
                       AND tm.chemin_media <> ''
                       AND tm.chemin_media <> '[]'
                       AND EXISTS (SELECT 1
                                   FROM jsonb_array_elements(tm.chemin_media::jsonb) AS elem
                                   WHERE COALESCE(elem ->> 'path', '') <> '')
                           ),
                       '[]'::jsonb
       )                     AS medias,
       -- Monuments lieux (IDs uniquement)
       COALESCE(array_agg(DISTINCT cml.monu_lieu_id) FILTER (WHERE cml.monu_lieu_id IS NOT NULL),
                '{}')        AS monuments_lieux_liees,
       -- Mobiliers images (IDs uniquement)
       COALESCE(array_agg(DISTINCT cpm.mobilier_image_id) FILTER (WHERE cpm.mobilier_image_id IS NOT NULL),
                '{}')        AS mobiliers_images_liees,
       -- Personnes morales (IDs uniquement)
       COALESCE(array_agg(DISTINCT cppmo.pers_morale_id) FILTER (WHERE cppmo.pers_morale_id IS NOT NULL),
                '{}')        AS personnes_morales_liees,
       -- Siècles
       COALESCE(array_agg(DISTINCT bs.siecle_list) FILTER (WHERE bs.siecle_list IS NOT NULL),
                '{}')        AS siecles,
       -- Themes
       COALESCE(array_agg(DISTINCT t.theme_type) FILTER (WHERE t.theme_type IS NOT NULL),
                '{}')        AS themes,
       publication_status,
       parent_id
FROM t_pers_physiques p
         LEFT JOIN cor_auteur_fiche_pers_phy cap ON p.id_pers_physique = cap.pers_physique_id
         LEFT JOIN bib_auteurs baf ON cap.auteur_fiche_pers_phy_id = baf.id_auteur_fiche
         LEFT JOIN loc_communes c ON p.id_commune = c.id_commune
         LEFT JOIN loc_departements d ON c.id_departement = d.id_departement
         LEFT JOIN loc_regions r ON d.id_region = r.id_region
         LEFT JOIN loc_pays pa ON r.id_pays = p.id_pays
         LEFT JOIN cor_medias_pers_phy cmp ON p.id_pers_physique = cmp.pers_physique_id
         LEFT JOIN t_medias tm ON cmp.media_pers_phy_id = tm.id_media
         LEFT JOIN cor_periodes_historiques_pers_phy cph ON p.id_pers_physique = cph.pers_physique_id
         LEFT JOIN bib_pers_phy_periodes_historiques bph ON cph.periode_historique_id = bph.id_periode_historique
         LEFT JOIN cor_modes_deplacements_pers_phy cmd ON p.id_pers_physique = cmd.pers_physique_id
         LEFT JOIN bib_pers_phy_modes_deplacements bmd ON cmd.mode_deplacement_id = bmd.id_mode_deplacement
         LEFT JOIN cor_professions_pers_phy cpp ON p.id_pers_physique = cpp.pers_physique_id
         LEFT JOIN bib_pers_phy_professions bpp ON cpp.profession_id = bpp.id_profession
         LEFT JOIN cor_monu_lieu_pers_phy cml ON p.id_pers_physique = cml.pers_phy_id
         LEFT JOIN cor_mob_img_pers_phy cpm ON p.id_pers_physique = cpm.pers_physique_id
         LEFT JOIN cor_pers_phy_pers_mo cppmo ON p.id_pers_physique = cppmo.pers_physique_id
         LEFT JOIN cor_siecles_pers_phy csp ON p.id_pers_physique = csp.pers_physique_id
         LEFT JOIN bib_siecle bs ON csp.siecle_pers_phy_id = bs.id_siecle
         LEFT JOIN cor_themes_pers_phy ctpp ON p.id_pers_physique = ctpp.pers_phy_id
         LEFT JOIN t_themes t ON t.id_theme = ctpp.theme_id
WHERE p.id_pers_physique = $1
GROUP BY p.id_pers_physique,
         c.nom_commune,
         d.nom_departement,
         r.nom_region,
         pa.nom_pays;

-- name: GetPendingPersonnesPhysiques :many
SELECT p.id_pers_physique     AS id,
       p.prenom_nom_pers_phy  AS firstname,
       p.date_naissance,
       p.date_deces,
       p.attestation,
       -- Periode historique
       COALESCE(array_agg(DISTINCT bph.periode_historique_type) FILTER (WHERE bph.periode_historique_type IS NOT NULL),
                '{}')         AS historical_period,
       p.bibliographie,
       p.elements_biographiques,
       p.elements_pelerinage,
       p.commutation_voeu,
       p.sources,
       p.date_creation,
       p.date_maj,
       p.publie,
       p.contributeurs,
       p.commentaires,
       -- Redacteurs (auteurs fiche)
       COALESCE(array_agg(DISTINCT baf.auteur_fiche_nom) FILTER (WHERE baf.auteur_fiche_nom IS NOT NULL),
                '{}')         AS redacteurs,
       -- Commune
       MAX(c.nom_commune)     AS commune,
       -- Département
       MAX(d.nom_departement) AS departement,
       -- Région
       MAX(r.nom_region)      AS region,
       -- Pays
       MAX(pa.nom_pays)       AS pays,
       -- Travels
       COALESCE(array_agg(DISTINCT bmd.mode_deplacement_type) FILTER (WHERE bmd.mode_deplacement_type IS NOT NULL),
                '{}')         AS travels,
       -- Professions
       COALESCE(array_agg(DISTINCT bpp.profession_type) FILTER (WHERE bpp.profession_type IS NOT NULL),
                '{}')         AS professions,
       p.nature_evenement,
       -- Médias
       COALESCE(
                       jsonb_agg(
                       DISTINCT jsonb_build_object(
                               'id', tm.id_media,
                               'titre', tm.titre_media
                                )
                                ) FILTER (
                           WHERE tm.chemin_media IS NOT NULL
                       AND tm.chemin_media <> ''
                       AND tm.chemin_media <> '[]'
                       AND EXISTS (SELECT 1
                                   FROM jsonb_array_elements(tm.chemin_media::jsonb) AS elem
                                   WHERE COALESCE(elem ->> 'path', '') <> '')
                           ),
                       '[]'::jsonb
       )                      AS medias,
       -- Monuments lieux (IDs uniquement)
       COALESCE(array_agg(DISTINCT cml.monu_lieu_id) FILTER (WHERE cml.monu_lieu_id IS NOT NULL),
                '{}')         AS monuments_lieux_liees,
       -- Mobiliers images (IDs uniquement)
       COALESCE(array_agg(DISTINCT cpm.mobilier_image_id) FILTER (WHERE cpm.mobilier_image_id IS NOT NULL),
                '{}')         AS mobiliers_images_liees,
       -- Personnes morales (IDs uniquement)
       COALESCE(array_agg(DISTINCT cppmo.pers_morale_id) FILTER (WHERE cppmo.pers_morale_id IS NOT NULL),
                '{}')         AS personnes_morales_liees,
       -- Siècles
       COALESCE(array_agg(DISTINCT bs.siecle_list) FILTER (WHERE bs.siecle_list IS NOT NULL),
                '{}')         AS siecles,
       -- Themes
       COALESCE(array_agg(DISTINCT t.theme_type) FILTER (WHERE t.theme_type IS NOT NULL),
                '{}')         AS themes
FROM t_pers_physiques p
         LEFT JOIN cor_auteur_fiche_pers_phy cap ON p.id_pers_physique = cap.pers_physique_id
         LEFT JOIN bib_auteurs baf ON cap.auteur_fiche_pers_phy_id = baf.id_auteur_fiche
         LEFT JOIN loc_communes c ON p.id_commune = c.id_commune
         LEFT JOIN loc_departements d ON c.id_departement = d.id_departement
         LEFT JOIN loc_regions r ON d.id_region = r.id_region
         LEFT JOIN loc_pays pa ON r.id_pays = p.id_pays
         LEFT JOIN cor_medias_pers_phy cmp ON p.id_pers_physique = cmp.pers_physique_id
         LEFT JOIN t_medias tm ON cmp.media_pers_phy_id = tm.id_media
         LEFT JOIN cor_periodes_historiques_pers_phy cph ON p.id_pers_physique = cph.pers_physique_id
         LEFT JOIN bib_pers_phy_periodes_historiques bph ON cph.periode_historique_id = bph.id_periode_historique
         LEFT JOIN cor_modes_deplacements_pers_phy cmd ON p.id_pers_physique = cmd.pers_physique_id
         LEFT JOIN bib_pers_phy_modes_deplacements bmd ON cmd.mode_deplacement_id = bmd.id_mode_deplacement
         LEFT JOIN cor_professions_pers_phy cpp ON p.id_pers_physique = cpp.pers_physique_id
         LEFT JOIN bib_pers_phy_professions bpp ON cpp.profession_id = bpp.id_profession
         LEFT JOIN cor_monu_lieu_pers_phy cml ON p.id_pers_physique = cml.pers_phy_id
         LEFT JOIN cor_mob_img_pers_phy cpm ON p.id_pers_physique = cpm.pers_physique_id
         LEFT JOIN cor_pers_phy_pers_mo cppmo ON p.id_pers_physique = cppmo.pers_physique_id
         LEFT JOIN cor_siecles_pers_phy csp ON p.id_pers_physique = csp.pers_physique_id
         LEFT JOIN bib_siecle bs ON csp.siecle_pers_phy_id = bs.id_siecle
         LEFT JOIN cor_themes_pers_phy ctpp ON p.id_pers_physique = ctpp.pers_phy_id
         LEFT JOIN t_themes t ON t.id_theme = ctpp.theme_id
WHERE p.publication_status = 'DRAFT'
   or p.publication_status = 'PENDING'
GROUP BY p.id_pers_physique;

-- name: ValidatePendingPersonnePhysique :exec
UPDATE t_pers_physiques
SET publication_status = 'PUBLISHED',
    publie             = true,
    parent_id          = NULL
WHERE id_pers_physique = $1;

-- name: DeletePendingPersonnePhysique :exec
DELETE
FROM t_pers_physiques
WHERE id_pers_physique = $1;

-- name: CreatePersPhysique :one
INSERT INTO t_pers_physiques
(prenom_nom_pers_phy,
 commentaires,
 date_naissance,
 date_deces,
 attestation,
 elements_biographiques,
 elements_pelerinage,
 nature_evenement,
 commutation_voeu,
 bibliographie,
 sources,
 date_creation,
 date_maj,
 contributeurs,
 id_commune,
 id_pays,
 publie,
 publication_status,
 parent_id)
VALUES (sqlc.arg(prenom_nom_pers_phy),
        sqlc.arg(commentaires),
        sqlc.arg(date_naissance),
        sqlc.arg(date_deces),
        sqlc.arg(attestation),
        sqlc.arg(elements_biographiques),
        sqlc.arg(elements_pelerinage),
        sqlc.arg(nature_evenement),
        sqlc.arg(commutation_voeu),
        sqlc.arg(bibliographie),
        sqlc.arg(sources),
        NOW(),
        NOW(),
        sqlc.arg(contributeurs),
        sqlc.arg(id_commune),
        sqlc.arg(id_pays),
        false,
        'DRAFT',
        sqlc.arg(parent_id))
RETURNING id_pers_physique;

-- name: AttachSieclesToPersPhy :exec
INSERT INTO cor_siecles_pers_phy
    (siecle_pers_phy_id, pers_physique_id)
SELECT unnest(sqlc.arg(siecle_id)::int[]), sqlc.arg(id);

-- name: AttachMediasToPersPhy :exec
INSERT INTO cor_medias_pers_phy
    (media_pers_phy_id, pers_physique_id)
SELECT unnest(sqlc.arg(media_ids)::int[]), sqlc.arg(id);

-- name: AttachThemesToPersPhy :exec
INSERT INTO cor_themes_pers_phy (theme_id, pers_phy_id)
SELECT unnest(sqlc.arg(theme_ids)::int[]), sqlc.arg(id);

-- name: AttachHistoricalPeriodsToPersPhy :exec
INSERT INTO cor_periodes_historiques_pers_phy (periode_historique_id, pers_physique_id)
SELECT unnest(sqlc.arg(periode_ids)::int[]), sqlc.arg(id);

-- name: AttachProfessionsToPersPhy :exec
INSERT INTO cor_professions_pers_phy (profession_id, pers_physique_id)
SELECT unnest(sqlc.arg(profession_ids)::int[]), sqlc.arg(id);

-- name: AttachModeDeTransportsToPersPhy :exec
INSERT INTO cor_modes_deplacements_pers_phy (mode_deplacement_id, pers_physique_id)
SELECT unnest(sqlc.arg(travel_ids)::int[]), sqlc.arg(id);

-- name: AttachAuthorToPersPhy :exec
INSERT INTO cor_auteur_fiche_pers_phy
    (auteur_fiche_pers_phy_id, pers_physique_id)
VALUES (sqlc.arg(auteur_id), sqlc.arg(id));

-- name: LinkPersPhyToMonuLieu :exec
INSERT INTO cor_monu_lieu_pers_phy
    (pers_phy_id, monu_lieu_id)
SELECT sqlc.arg(id), unnest(sqlc.arg(monu_lieu_ids)::int[]);

-- name: LinkPersPhyToMobImg :exec
INSERT INTO cor_mob_img_pers_phy
    (pers_physique_id, mobilier_image_id)
SELECT sqlc.arg(id), unnest(sqlc.arg(mob_img_ids)::int[]);

-- name: LinkPersPhyToPersMo :exec
INSERT INTO cor_pers_phy_pers_mo
    (pers_physique_id, pers_morale_id)
SELECT sqlc.arg(id), unnest(sqlc.arg(perso_mo_ids)::int[]);

-- name: DetachSieclesFromPersPhy :exec
DELETE
FROM cor_siecles_pers_phy
WHERE pers_physique_id = sqlc.arg(id);

-- name: DetachMediasFromPersPhy :exec
DELETE
FROM cor_medias_pers_phy
WHERE pers_physique_id = sqlc.arg(id);

-- name: DetachThemesFromPersPhy :exec
DELETE
FROM cor_themes_pers_phy
WHERE pers_phy_id = sqlc.arg(id);

-- name: DetachHistoricalPeriodsFromPersPhy :exec
DELETE
FROM cor_periodes_historiques_pers_phy
WHERE pers_physique_id = sqlc.arg(id);

-- name: DetachProfessionsFromPersPhy :exec
DELETE
FROM cor_professions_pers_phy
WHERE pers_physique_id = sqlc.arg(id);

-- name: DetachModeDeTransportsFromPersPhy :exec
DELETE
FROM cor_modes_deplacements_pers_phy
WHERE pers_physique_id = sqlc.arg(id);

-- name: DetachAuthorFromPersPhy :exec
DELETE
FROM cor_auteur_fiche_pers_phy
WHERE pers_physique_id = sqlc.arg(id);

-- name: UnlinkPersPhyFromMonuLieu :exec
DELETE
FROM cor_monu_lieu_pers_phy
WHERE pers_phy_id = sqlc.arg(id);

-- name: UnlinkPersPhyFromMobImg :exec
DELETE
FROM cor_mob_img_pers_phy
WHERE pers_physique_id = sqlc.arg(id);

-- name: UnlinkPersPhyFromPersMo :exec
DELETE
FROM cor_pers_phy_pers_mo
WHERE pers_physique_id = sqlc.arg(id);