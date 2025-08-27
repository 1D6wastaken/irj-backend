-- name: GetFilteredMonumentsLieux :many
SELECT m.id_monument_lieu                                                                             AS id,
       m.titre_monu_lieu                                                                              AS title,
       COALESCE(array_agg(DISTINCT bs.siecle_list) FILTER (WHERE bs.siecle_list IS NOT NULL), '{}')   AS siecles,
       COALESCE(array_agg(DISTINCT bmn.monu_lieu_nature_type) FILTER (WHERE bmn.monu_lieu_nature_type IS NOT NULL),
                '{}')                                                                                 AS natures,
       COALESCE(array_agg(DISTINCT tm.chemin_media) FILTER (WHERE tm.chemin_media IS NOT NULL), '{}') AS medias
FROM t_monuments_lieux m
         LEFT JOIN cor_siecles_monu_lieu csl ON m.id_monument_lieu = csl.monument_lieu_id
         LEFT JOIN bib_siecle bs ON bs.id_siecle = csl.siecle_monu_lieu_id

         LEFT JOIN cor_natures_monu_lieu cnl ON m.id_monument_lieu = cnl.monument_lieu_id
         LEFT JOIN bib_monu_lieu_natures bmn ON bmn.id_monu_lieu_nature = cnl.monu_lieu_nature_id

         LEFT JOIN cor_medias_monu_lieu cme ON m.id_monument_lieu = cme.monument_lieu_id
         LEFT JOIN t_medias tm ON tm.id_media = cme.media_monu_lieu_id

         LEFT JOIN loc_communes lc ON m.id_commune = lc.id_commune
         LEFT JOIN loc_departements ld ON lc.id_departement = ld.id_departement
         LEFT JOIN loc_regions lr ON ld.id_region = lr.id_region
         LEFT JOIN loc_pays lp ON lr.id_pays = lp.id_pays

         LEFT JOIN cor_etat_cons_monu_lieu cel ON m.id_monument_lieu = cel.monument_lieu_id
         LEFT JOIN cor_auteur_fiche_monu_lieu caf ON m.id_monument_lieu = caf.monument_lieu_id
         LEFT JOIN cor_materiaux_monu_lieu cml ON m.id_monument_lieu = cml.monument_lieu_id
         LEFT JOIN cor_monu_lieu_pers_mo cpm ON m.id_monument_lieu = cpm.monument_lieu_id
         LEFT JOIN cor_monu_lieu_pers_phy cpp ON m.id_monument_lieu = cpp.monu_lieu_id
         LEFT JOIN cor_monu_lieu_mob_img cmi ON m.id_monument_lieu = cmi.monument_lieu_id
WHERE (cardinality(sqlc.arg(siecles)::int[]) = 0 OR csl.siecle_monu_lieu_id = ANY (sqlc.arg(siecles)::int[]))
   OR (cardinality(sqlc.arg(natures)::int[]) = 0 OR cnl.monu_lieu_nature_id = ANY (sqlc.arg(natures)::int[]))
   OR (cardinality(sqlc.arg(etats)::int[]) = 0 OR cel.etat_cons_monu_lieu_id = ANY (sqlc.arg(etats)::int[]))
   OR (cardinality(sqlc.arg(auteurs_fiche)::int[]) = 0 OR
       caf.auteur_fiche_monu_lieu_id = ANY (sqlc.arg(auteurs_fiche)::int[]))
   OR (cardinality(sqlc.arg(materiaux)::int[]) = 0 OR cml.materiau_monu_lieu_id = ANY (sqlc.arg(materiaux)::int[]))
   OR (cardinality(sqlc.arg(pers_mo)::int[]) = 0 OR cpm.pers_morale_id = ANY (sqlc.arg(pers_mo)::int[]))
   OR (cardinality(sqlc.arg(pers_phy)::int[]) = 0 OR cpp.pers_phy_id = ANY (sqlc.arg(pers_phy)::int[]))
   OR (cardinality(sqlc.arg(mobiliers)::int[]) = 0 OR cmi.mobilier_image_id = ANY (sqlc.arg(mobiliers)::int[]))
   OR (cardinality(sqlc.arg(cities)::int[]) = 0 OR lc.id_commune = ANY (sqlc.arg(cities)::int[]))
   OR (cardinality(sqlc.arg(departments)::int[]) = 0 OR ld.id_departement = ANY (sqlc.arg(departments)::int[]))
   OR (cardinality(sqlc.arg(regions)::int[]) = 0 OR lr.id_region = ANY (sqlc.arg(regions)::int[]))
   OR (cardinality(sqlc.arg(pays)::int[]) = 0 OR lp.id_pays = ANY (sqlc.arg(pays)::int[]))
GROUP BY m.id_monument_lieu
ORDER BY m.id_monument_lieu
LIMIT sqlc.arg(limit_param) OFFSET sqlc.arg(offset_param);

-- name: GetMonumentLieuByID :one
SELECT m.id_monument_lieu AS id,
       m.titre_monu_lieu  AS title,
       m.description,
       m.histoire,
       m.geolocalisation,
       m.bibliographie,
       m.date_creation,
       m.date_maj,
       m.publie,
       m.contributeurs,
       m.protection,
       m.protection_commentaires,
       m.source,
       -- Redacteurs (auteurs fiche)
       COALESCE(array_agg(DISTINCT baf.auteur_fiche_nom) FILTER (WHERE baf.auteur_fiche_nom IS NOT NULL),
                '{}')     AS redacteurs,
       -- Commune
       c.nom_commune      AS commune,
       -- Département
       d.nom_departement  AS departement,
       -- Région
       r.nom_region       AS region,
       -- Pays
       p.nom_pays         AS pays,
       -- États de conservation
       COALESCE(array_agg(DISTINCT bec.etat_conservation_type) FILTER (WHERE bec.etat_conservation_type IS NOT NULL),
                '{}')     AS etats_conservation,
       -- Matériaux
       COALESCE(array_agg(DISTINCT bm.materiau_type) FILTER (WHERE bm.materiau_type IS NOT NULL),
                '{}')     AS materiaux,
       -- Natures
       COALESCE(array_agg(DISTINCT bmn.monu_lieu_nature_type) FILTER (WHERE bmn.id_monu_lieu_nature IS NOT NULL),
                '{}')     AS natures,
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
       )                  AS medias,
       -- Mobiliers images (IDs uniquement)
       COALESCE(array_agg(DISTINCT cmi.mobilier_image_id) FILTER (WHERE cmi.mobilier_image_id IS NOT NULL),
                '{}')     AS mobiliers_images_liees,
       -- Personnes morales (IDs uniquement)
       COALESCE(array_agg(DISTINCT cpm.pers_morale_id) FILTER (WHERE cpm.pers_morale_id IS NOT NULL),
                '{}')     AS personnes_morales_liees,
       -- Personnes physiques (IDs uniquement)
       COALESCE(array_agg(DISTINCT cpp.pers_phy_id) FILTER (WHERE cpp.pers_phy_id IS NOT NULL),
                '{}')     AS personnes_physiques_liees,
       -- Siècles
       COALESCE(array_agg(DISTINCT bs.siecle_list) FILTER (WHERE bs.siecle_list IS NOT NULL),
                '{}')     AS siecles,
       -- Themes
       COALESCE(array_agg(DISTINCT t.theme_type) FILTER (WHERE t.theme_type IS NOT NULL),
                '{}')     AS themes,
       publication_status,
       parent_id
FROM t_monuments_lieux m
         LEFT JOIN cor_auteur_fiche_monu_lieu caf ON m.id_monument_lieu = caf.monument_lieu_id
         LEFT JOIN bib_auteurs baf ON caf.auteur_fiche_monu_lieu_id = baf.id_auteur_fiche
         LEFT JOIN loc_communes c ON m.id_commune = c.id_commune
         LEFT JOIN loc_departements d ON c.id_departement = d.id_departement
         LEFT JOIN loc_regions r ON d.id_region = r.id_region
         LEFT JOIN loc_pays p ON r.id_pays = p.id_pays
         LEFT JOIN cor_etat_cons_monu_lieu cec ON m.id_monument_lieu = cec.monument_lieu_id
         LEFT JOIN bib_etats_conservation bec ON cec.etat_cons_monu_lieu_id = bec.id_etat_conservation
         LEFT JOIN cor_materiaux_monu_lieu cm ON m.id_monument_lieu = cm.monument_lieu_id
         LEFT JOIN bib_materiaux bm ON cm.materiau_monu_lieu_id = bm.id_materiau
         LEFT JOIN cor_medias_monu_lieu cmm ON m.id_monument_lieu = cmm.monument_lieu_id
         LEFT JOIN cor_natures_monu_lieu cnm ON m.id_monument_lieu = cnm.monument_lieu_id
         LEFT JOIN bib_monu_lieu_natures bmn ON cnm.monu_lieu_nature_id = bmn.id_monu_lieu_nature
         LEFT JOIN t_medias tm ON cmm.media_monu_lieu_id = tm.id_media
         LEFT JOIN cor_monu_lieu_mob_img cmi ON m.id_monument_lieu = cmi.monument_lieu_id
         LEFT JOIN cor_monu_lieu_pers_mo cpm ON m.id_monument_lieu = cpm.monument_lieu_id
         LEFT JOIN cor_monu_lieu_pers_phy cpp ON m.id_monument_lieu = cpp.monu_lieu_id
         LEFT JOIN cor_siecles_monu_lieu csl ON m.id_monument_lieu = csl.monument_lieu_id
         LEFT JOIN bib_siecle bs ON csl.siecle_monu_lieu_id = bs.id_siecle
         LEFT JOIN cor_themes_monu_lieu ctml ON m.id_monument_lieu = ctml.monu_lieu_id
         LEFT JOIN t_themes t ON t.id_theme = ctml.theme_id
WHERE m.id_monument_lieu = $1
GROUP BY m.id_monument_lieu,
         c.nom_commune,
         d.nom_departement,
         r.nom_region,
         p.nom_pays;

-- name: GetPendingMonumentsLieux :many
SELECT m.id_monument_lieu     AS id,
       m.titre_monu_lieu      AS title,
       m.description,
       m.histoire,
       m.geolocalisation,
       m.bibliographie,
       m.date_creation,
       m.date_maj,
       m.publie,
       m.contributeurs,
       m.protection,
       m.protection_commentaires,
       m.source,
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
       MAX(p.nom_pays)        AS pays,
       -- États de conservation
       COALESCE(array_agg(DISTINCT bec.etat_conservation_type) FILTER (WHERE bec.etat_conservation_type IS NOT NULL),
                '{}')         AS etats_conservation,
       -- Matériaux
       COALESCE(array_agg(DISTINCT bm.materiau_type) FILTER (WHERE bm.materiau_type IS NOT NULL),
                '{}')         AS materiaux,
       -- Natures
       COALESCE(array_agg(DISTINCT bmn.monu_lieu_nature_type) FILTER (WHERE bmn.id_monu_lieu_nature IS NOT NULL),
                '{}')         AS natures,
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
       -- Mobiliers images (IDs uniquement)
       COALESCE(array_agg(DISTINCT cmi.mobilier_image_id) FILTER (WHERE cmi.mobilier_image_id IS NOT NULL),
                '{}')         AS mobiliers_images_liees,
       -- Personnes morales (IDs uniquement)
       COALESCE(array_agg(DISTINCT cpm.pers_morale_id) FILTER (WHERE cpm.pers_morale_id IS NOT NULL),
                '{}')         AS personnes_morales_liees,
       -- Personnes physiques (IDs uniquement)
       COALESCE(array_agg(DISTINCT cpp.pers_phy_id) FILTER (WHERE cpp.pers_phy_id IS NOT NULL),
                '{}')         AS personnes_physiques_liees,
       -- Siècles
       COALESCE(array_agg(DISTINCT bs.siecle_list) FILTER (WHERE bs.siecle_list IS NOT NULL),
                '{}')         AS siecles,
       -- Themes
       COALESCE(array_agg(DISTINCT t.theme_type) FILTER (WHERE t.theme_type IS NOT NULL),
                '{}')         AS themes
FROM t_monuments_lieux m
         LEFT JOIN cor_auteur_fiche_monu_lieu caf ON m.id_monument_lieu = caf.monument_lieu_id
         LEFT JOIN bib_auteurs baf ON caf.auteur_fiche_monu_lieu_id = baf.id_auteur_fiche
         LEFT JOIN loc_communes c ON m.id_commune = c.id_commune
         LEFT JOIN loc_departements d ON c.id_departement = d.id_departement
         LEFT JOIN loc_regions r ON d.id_region = r.id_region
         LEFT JOIN loc_pays p ON r.id_pays = p.id_pays
         LEFT JOIN cor_etat_cons_monu_lieu cec ON m.id_monument_lieu = cec.monument_lieu_id
         LEFT JOIN bib_etats_conservation bec ON cec.etat_cons_monu_lieu_id = bec.id_etat_conservation
         LEFT JOIN cor_materiaux_monu_lieu cm ON m.id_monument_lieu = cm.monument_lieu_id
         LEFT JOIN bib_materiaux bm ON cm.materiau_monu_lieu_id = bm.id_materiau
         LEFT JOIN cor_medias_monu_lieu cmm ON m.id_monument_lieu = cmm.monument_lieu_id
         LEFT JOIN cor_natures_monu_lieu cnm ON m.id_monument_lieu = cnm.monument_lieu_id
         LEFT JOIN bib_monu_lieu_natures bmn ON cnm.monu_lieu_nature_id = bmn.id_monu_lieu_nature
         LEFT JOIN t_medias tm ON cmm.media_monu_lieu_id = tm.id_media
         LEFT JOIN cor_monu_lieu_mob_img cmi ON m.id_monument_lieu = cmi.monument_lieu_id
         LEFT JOIN cor_monu_lieu_pers_mo cpm ON m.id_monument_lieu = cpm.monument_lieu_id
         LEFT JOIN cor_monu_lieu_pers_phy cpp ON m.id_monument_lieu = cpp.monu_lieu_id
         LEFT JOIN cor_siecles_monu_lieu csl ON m.id_monument_lieu = csl.monument_lieu_id
         LEFT JOIN bib_siecle bs ON csl.siecle_monu_lieu_id = bs.id_siecle
         LEFT JOIN cor_themes_monu_lieu ctml ON m.id_monument_lieu = ctml.monu_lieu_id
         LEFT JOIN t_themes t ON t.id_theme = ctml.theme_id
WHERE m.publication_status = 'DRAFT'
   OR m.publication_status = 'PENDING'
GROUP BY m.id_monument_lieu;

-- name: ValidatePendingMonumentLieu :exec
UPDATE t_monuments_lieux
SET publication_status = 'PUBLISHED',
    publie             = true,
    parent_id          = NULL
WHERE id_monument_lieu = $1;

-- name: DeletePendingMonumentLieu :exec
DELETE
FROM t_monuments_lieux
WHERE id_monument_lieu = $1;

-- name: CreateMonumentLieu :one
INSERT INTO t_monuments_lieux
(titre_monu_lieu,
 description,
 histoire,
 geolocalisation,
 bibliographie,
 protection,
 protection_commentaires,
 source,
 date_creation,
 date_maj,
 contributeurs,
 id_commune,
 id_pays,
 publie,
 publication_status,
 parent_id)
VALUES (sqlc.arg(titre_monu_lieu),
        sqlc.arg(description),
        sqlc.arg(histoire),
        sqlc.arg(geolocalisation),
        sqlc.arg(bibliographie),
        sqlc.arg(protection),
        sqlc.arg(protection_commentaires),
        sqlc.arg(source),
        NOW(),
        NOW(),
        sqlc.arg(contributeurs),
        sqlc.arg(id_commune),
        sqlc.arg(id_pays),
        false,
        'DRAFT',
        sqlc.arg(parent_id))
RETURNING id_monument_lieu;

-- name: AttachSieclesToMonuLieu :exec
INSERT INTO cor_siecles_monu_lieu
    (siecle_monu_lieu_id, monument_lieu_id)
SELECT unnest(sqlc.arg(siecle_id)::int[]), sqlc.arg(id);

-- name: AttachMediasToMonuLieu :exec
INSERT INTO cor_medias_monu_lieu
    (media_monu_lieu_id, monument_lieu_id)
SELECT unnest(sqlc.arg(media_ids)::int[]), sqlc.arg(id);

-- name: AttachThemesToMonuLieu :exec
INSERT INTO cor_themes_monu_lieu (theme_id, monu_lieu_id)
SELECT unnest(sqlc.arg(theme_ids)::int[]), sqlc.arg(id);

-- name: AttachNaturesToMonuLieu :exec
INSERT INTO cor_natures_monu_lieu (monu_lieu_nature_id, monument_lieu_id)
SELECT unnest(sqlc.arg(nature_ids)::int[]), sqlc.arg(id);

-- name: AttachEtatsToMonuLieu :exec
INSERT INTO cor_etat_cons_monu_lieu
    (etat_cons_monu_lieu_id, monument_lieu_id)
SELECT unnest(sqlc.arg(etat_ids)::int[]), sqlc.arg(id);

-- name: AttachAuthorToMonuLieu :exec
INSERT INTO cor_auteur_fiche_monu_lieu
    (auteur_fiche_monu_lieu_id, monument_lieu_id)
VALUES (sqlc.arg(auteur_id), sqlc.arg(id));

-- name: AttachMateriauxToMonuLieu :exec
INSERT INTO cor_materiaux_monu_lieu
    (materiau_monu_lieu_id, monument_lieu_id)
SELECT unnest(sqlc.arg(materiau_ids)::int[]), sqlc.arg(id);

-- name: LinkMonuLieuToMobImg :exec
INSERT INTO cor_monu_lieu_mob_img
    (monument_lieu_id, mobilier_image_id)
SELECT sqlc.arg(id), unnest(sqlc.arg(mob_img_ids)::int[]);

-- name: LinkMonuLieuToPersMo :exec
INSERT INTO cor_monu_lieu_pers_mo
    (monument_lieu_id, pers_morale_id)
SELECT sqlc.arg(id), unnest(sqlc.arg(perso_mo_ids)::int[]);

-- name: LinkMonuLieuToPersPhy :exec
INSERT INTO cor_monu_lieu_pers_phy
    (monu_lieu_id, pers_phy_id)
SELECT sqlc.arg(id), unnest(sqlc.arg(perso_phy_ids)::int[]);

-- name: DetachSieclesFromMonuLieu :exec
DELETE
FROM cor_siecles_monu_lieu
WHERE monument_lieu_id = sqlc.arg(id);

-- name: DetachMediasFromMonuLieu :exec
DELETE
FROM cor_medias_monu_lieu
WHERE monument_lieu_id = sqlc.arg(id);

-- name: DetachThemesFromMonuLieu :exec
DELETE
FROM cor_themes_monu_lieu
WHERE monu_lieu_id = sqlc.arg(id);

-- name: DetachNaturesFromMonuLieu :exec
DELETE
FROM cor_natures_monu_lieu
WHERE monument_lieu_id = sqlc.arg(id);

-- name: DetachEtatsFromMonuLieu :exec
DELETE
FROM cor_etat_cons_monu_lieu
WHERE monument_lieu_id = sqlc.arg(id);

-- name: DetachAuthorFromMonuLieu :exec
DELETE
FROM cor_auteur_fiche_monu_lieu
WHERE monument_lieu_id = sqlc.arg(id);

-- name: DetachMateriauxFromMonuLieu :exec
DELETE
FROM cor_materiaux_monu_lieu
WHERE monument_lieu_id = sqlc.arg(id);

-- name: UnlinkMonuLieuFromMobImg :exec
DELETE
FROM cor_monu_lieu_mob_img
WHERE monument_lieu_id = sqlc.arg(id);

-- name: UnlinkMonuLieuFromPersMo :exec
DELETE
FROM cor_monu_lieu_pers_mo
WHERE monument_lieu_id = sqlc.arg(id);

-- name: UnlinkMonuLieuFromPersPhy :exec
DELETE
FROM cor_monu_lieu_pers_phy
WHERE monu_lieu_id = sqlc.arg(id);