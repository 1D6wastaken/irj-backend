-- name: GetFilteredMobiliersImages :many
SELECT m.id_mobilier_image                                                                            AS id,
       m.titre_mob_img                                                                                AS title,
       COALESCE(array_agg(DISTINCT bs.siecle_list) FILTER (WHERE bs.siecle_list IS NOT NULL), '{}')   AS siecles,
       COALESCE(array_agg(DISTINCT bmn.nature_type) FILTER (WHERE bmn.nature_type IS NOT NULL), '{}') AS natures,
       COALESCE(array_agg(DISTINCT tm.chemin_media) FILTER (WHERE tm.chemin_media IS NOT NULL), '{}') AS medias
FROM t_mobiliers_images m
         LEFT JOIN cor_siecles_mob_img csm ON m.id_mobilier_image = csm.mobilier_image_id
         LEFT JOIN bib_siecle bs ON bs.id_siecle = csm.siecle_mob_img_id

         LEFT JOIN cor_natures_mob_img cnm ON m.id_mobilier_image = cnm.mobilier_image_id
         LEFT JOIN bib_mob_img_natures bmn ON bmn.id_nature = cnm.nature_id

         LEFT JOIN cor_techniques_mob_img ctm ON m.id_mobilier_image = ctm.mobilier_image_id
         LEFT JOIN bib_mob_img_techniques bmt ON bmt.id_technique = ctm.technique_id

         LEFT JOIN cor_medias_mob_img cme ON m.id_mobilier_image = cme.mobilier_image_id
         LEFT JOIN t_medias tm ON tm.id_media = cme.media_mob_img_id

         LEFT JOIN loc_communes c ON m.id_commune = c.id_commune
         LEFT JOIN loc_departements d ON c.id_departement = d.id_departement
         LEFT JOIN loc_regions r ON d.id_region = r.id_region
         LEFT JOIN loc_pays p ON r.id_pays = p.id_pays

         LEFT JOIN cor_etat_cons_mob_img cel ON m.id_mobilier_image = cel.mobilier_image_id
         LEFT JOIN cor_auteur_fiche_mob_img caf ON m.id_mobilier_image = caf.mobilier_image_id
         LEFT JOIN cor_materiaux_mob_img cml ON m.id_mobilier_image = cml.mobilier_image_id
         LEFT JOIN cor_mob_img_pers_mo cpm ON m.id_mobilier_image = cpm.mobilier_image_id
         LEFT JOIN cor_mob_img_pers_phy cpp ON m.id_mobilier_image = cpp.mobilier_image_id
         LEFT JOIN cor_monu_lieu_mob_img cmi ON m.id_mobilier_image = cmi.mobilier_image_id
WHERE (cardinality(sqlc.arg(siecles)::int[]) = 0 OR csm.siecle_mob_img_id = ANY (sqlc.arg(siecles)::int[]))
   OR (cardinality(sqlc.arg(natures)::int[]) = 0 OR cnm.nature_id = ANY (sqlc.arg(natures)::int[]))
   OR (cardinality(sqlc.arg(techniques)::int[]) = 0 OR ctm.technique_id = ANY (sqlc.arg(techniques)::int[]))
   OR (cardinality(sqlc.arg(etats)::int[]) = 0 OR cel.etat_cons_mob_img_id = ANY (sqlc.arg(etats)::int[]))
   OR (cardinality(sqlc.arg(auteurs_fiche)::int[]) = 0 OR
       caf.auteur_fiche_mob_img_id = ANY (sqlc.arg(auteurs_fiche)::int[]))
   OR (cardinality(sqlc.arg(materiaux)::int[]) = 0 OR cml.materiau_mob_img_id = ANY (sqlc.arg(materiaux)::int[]))
   OR (cardinality(sqlc.arg(pers_mo)::int[]) = 0 OR cpm.pers_morale_id = ANY (sqlc.arg(pers_mo)::int[]))
   OR (cardinality(sqlc.arg(pers_phy)::int[]) = 0 OR cpp.pers_physique_id = ANY (sqlc.arg(pers_phy)::int[]))
   OR (cardinality(sqlc.arg(places)::int[]) = 0 OR cmi.monument_lieu_id = ANY (sqlc.arg(places)::int[]))
   OR (cardinality(sqlc.arg(cities)::int[]) = 0 OR c.id_commune = ANY (sqlc.arg(cities)::int[]))
   OR (cardinality(sqlc.arg(departments)::int[]) = 0 OR d.id_departement = ANY (sqlc.arg(departments)::int[]))
   OR (cardinality(sqlc.arg(regions)::int[]) = 0 OR r.id_region = ANY (sqlc.arg(regions)::int[]))
   OR (cardinality(sqlc.arg(pays)::int[]) = 0 OR p.id_pays = ANY (sqlc.arg(pays)::int[]))
GROUP BY m.id_mobilier_image
ORDER BY m.id_mobilier_image
LIMIT sqlc.arg(limit_param) OFFSET sqlc.arg(offset_param);

-- name: GetMobilierImageByID :one
SELECT m.id_mobilier_image AS id,
       m.titre_mob_img     AS title,
       m.description,
       m.historique,
       m.bibliographie,
       m.inscriptions,
       m.date_cr_ation,
       m.date_maj,
       m.publie,
       m.contributeurs,
       m.source,
       m.protection,
       m.protection_commentaires,
       m.lieu_conservation,
       m.lieu_origine,
       -- Redacteurs (auteurs fiche)
       COALESCE(array_agg(DISTINCT baf.auteur_fiche_nom) FILTER (WHERE baf.auteur_fiche_nom IS NOT NULL),
                '{}')      AS redacteurs,
       -- Commune
       c.nom_commune       AS commune,
       -- Département
       d.nom_departement   AS departement,
       -- Région
       r.nom_region        AS region,
       -- Pays
       p.nom_pays          AS pays,
       -- États de conservation
       COALESCE(array_agg(DISTINCT bec.etat_conservation_type) FILTER (WHERE bec.etat_conservation_type IS NOT NULL),
                '{}')      AS etats_conservation,
       -- Matériaux
       COALESCE(array_agg(DISTINCT bm.materiau_type) FILTER (WHERE bm.materiau_type IS NOT NULL),
                '{}')      AS materiaux,
       -- Techniques
       COALESCE(array_agg(DISTINCT bmt.technique_type) FILTER (WHERE bmt.technique_type IS NOT NULL),
                '{}')      AS techniques,
       -- Natures
       COALESCE(array_agg(DISTINCT bmn.nature_type) FILTER (WHERE bmn.nature_type IS NOT NULL),
                '{}')      AS natures,
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
       )                   AS medias,
       -- Monuments lieux (IDs uniquement)
       COALESCE(array_agg(DISTINCT cmi.monument_lieu_id) FILTER (WHERE cmi.monument_lieu_id IS NOT NULL),
                '{}')      AS monuments_lieux_liees,
       -- Personnes morales (IDs uniquement)
       COALESCE(array_agg(DISTINCT cpm.pers_morale_id) FILTER (WHERE cpm.pers_morale_id IS NOT NULL),
                '{}')      AS personnes_morales_liees,
       -- Personnes physiques (IDs uniquement)
       COALESCE(array_agg(DISTINCT cpp.pers_physique_id) FILTER (WHERE cpp.pers_physique_id IS NOT NULL),
                '{}')      AS personnes_physiques_liees,
       -- Siècles
       COALESCE(array_agg(DISTINCT bs.siecle_list) FILTER (WHERE bs.siecle_list IS NOT NULL),
                '{}')      AS siecles,
       -- Themes
       COALESCE(array_agg(DISTINCT t.theme_type) FILTER (WHERE t.theme_type IS NOT NULL),
                '{}')      AS themes,
       publication_status,
       parent_id
FROM t_mobiliers_images m
         LEFT JOIN cor_auteur_fiche_mob_img caf ON m.id_mobilier_image = caf.mobilier_image_id
         LEFT JOIN bib_auteurs baf ON caf.auteur_fiche_mob_img_id = baf.id_auteur_fiche
         LEFT JOIN loc_communes c ON m.id_commune = c.id_commune
         LEFT JOIN loc_departements d ON c.id_departement = d.id_departement
         LEFT JOIN loc_regions r ON d.id_region = r.id_region
         LEFT JOIN loc_pays p ON r.id_pays = p.id_pays
         LEFT JOIN cor_techniques_mob_img ctm ON m.id_mobilier_image = ctm.mobilier_image_id
         LEFT JOIN bib_mob_img_techniques bmt ON bmt.id_technique = ctm.technique_id
         LEFT JOIN cor_etat_cons_mob_img cec ON m.id_mobilier_image = cec.mobilier_image_id
         LEFT JOIN bib_etats_conservation bec ON cec.etat_cons_mob_img_id = bec.id_etat_conservation
         LEFT JOIN cor_materiaux_mob_img cm ON m.id_mobilier_image = cm.mobilier_image_id
         LEFT JOIN bib_materiaux bm ON cm.materiau_mob_img_id = bm.id_materiau
         LEFT JOIN cor_medias_mob_img cmm ON m.id_mobilier_image = cmm.mobilier_image_id
         LEFT JOIN cor_natures_mob_img cnm ON m.id_mobilier_image = cnm.mobilier_image_id
         LEFT JOIN bib_mob_img_natures bmn ON cnm.nature_id = bmn.id_nature
         LEFT JOIN t_medias tm ON cmm.media_mob_img_id = tm.id_media
         LEFT JOIN cor_monu_lieu_mob_img cmi ON m.id_mobilier_image = cmi.mobilier_image_id
         LEFT JOIN cor_mob_img_pers_mo cpm ON m.id_mobilier_image = cpm.mobilier_image_id
         LEFT JOIN cor_mob_img_pers_phy cpp ON m.id_mobilier_image = cpp.mobilier_image_id
         LEFT JOIN cor_siecles_mob_img csl ON m.id_mobilier_image = csl.mobilier_image_id
         LEFT JOIN bib_siecle bs ON csl.siecle_mob_img_id = bs.id_siecle
         LEFT JOIN cor_themes_mob_img ctmi ON m.id_mobilier_image = ctmi.mob_img_id
         LEFT JOIN t_themes t ON t.id_theme = ctmi.theme_id
WHERE m.id_mobilier_image = $1
GROUP BY m.id_mobilier_image,
         c.nom_commune,
         d.nom_departement,
         r.nom_region,
         p.nom_pays;

-- name: GetPendingMobiliersImages :many
SELECT m.id_mobilier_image    AS id,
       m.titre_mob_img        AS title,
       m.description,
       m.historique,
       m.bibliographie,
       m.inscriptions,
       m.date_cr_ation,
       m.date_maj,
       m.publie,
       m.contributeurs,
       m.source,
       m.protection,
       m.protection_commentaires,
       m.lieu_conservation,
       m.lieu_origine,
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
       -- Techniques
       COALESCE(array_agg(DISTINCT bmt.technique_type) FILTER (WHERE bmt.technique_type IS NOT NULL),
                '{}')         AS techniques,
       -- Natures
       COALESCE(array_agg(DISTINCT bmn.nature_type) FILTER (WHERE bmn.nature_type IS NOT NULL),
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
       -- Monuments lieux (IDs uniquement)
       COALESCE(array_agg(DISTINCT cmi.monument_lieu_id) FILTER (WHERE cmi.monument_lieu_id IS NOT NULL),
                '{}')         AS monuments_lieux_liees,
       -- Personnes morales (IDs uniquement)
       COALESCE(array_agg(DISTINCT cpm.pers_morale_id) FILTER (WHERE cpm.pers_morale_id IS NOT NULL),
                '{}')         AS personnes_morales_liees,
       -- Personnes physiques (IDs uniquement)
       COALESCE(array_agg(DISTINCT cpp.pers_physique_id) FILTER (WHERE cpp.pers_physique_id IS NOT NULL),
                '{}')         AS personnes_physiques_liees,
       -- Siècles
       COALESCE(array_agg(DISTINCT bs.siecle_list) FILTER (WHERE bs.siecle_list IS NOT NULL),
                '{}')         AS siecles,
       -- Themes
       COALESCE(array_agg(DISTINCT t.theme_type) FILTER (WHERE t.theme_type IS NOT NULL),
                '{}')         AS themes
FROM t_mobiliers_images m
         LEFT JOIN cor_auteur_fiche_mob_img caf ON m.id_mobilier_image = caf.mobilier_image_id
         LEFT JOIN bib_auteurs baf ON caf.auteur_fiche_mob_img_id = baf.id_auteur_fiche
         LEFT JOIN loc_communes c ON m.id_commune = c.id_commune
         LEFT JOIN loc_departements d ON c.id_departement = d.id_departement
         LEFT JOIN loc_regions r ON d.id_region = r.id_region
         LEFT JOIN loc_pays p ON r.id_pays = p.id_pays
         LEFT JOIN cor_techniques_mob_img ctm ON m.id_mobilier_image = ctm.mobilier_image_id
         LEFT JOIN bib_mob_img_techniques bmt ON bmt.id_technique = ctm.technique_id
         LEFT JOIN cor_etat_cons_mob_img cec ON m.id_mobilier_image = cec.mobilier_image_id
         LEFT JOIN bib_etats_conservation bec ON cec.etat_cons_mob_img_id = bec.id_etat_conservation
         LEFT JOIN cor_materiaux_mob_img cm ON m.id_mobilier_image = cm.mobilier_image_id
         LEFT JOIN bib_materiaux bm ON cm.materiau_mob_img_id = bm.id_materiau
         LEFT JOIN cor_medias_mob_img cmm ON m.id_mobilier_image = cmm.mobilier_image_id
         LEFT JOIN cor_natures_mob_img cnm ON m.id_mobilier_image = cnm.mobilier_image_id
         LEFT JOIN bib_mob_img_natures bmn ON cnm.nature_id = bmn.id_nature
         LEFT JOIN t_medias tm ON cmm.media_mob_img_id = tm.id_media
         LEFT JOIN cor_monu_lieu_mob_img cmi ON m.id_mobilier_image = cmi.mobilier_image_id
         LEFT JOIN cor_mob_img_pers_mo cpm ON m.id_mobilier_image = cpm.mobilier_image_id
         LEFT JOIN cor_mob_img_pers_phy cpp ON m.id_mobilier_image = cpp.mobilier_image_id
         LEFT JOIN cor_siecles_mob_img csl ON m.id_mobilier_image = csl.mobilier_image_id
         LEFT JOIN bib_siecle bs ON csl.siecle_mob_img_id = bs.id_siecle
         LEFT JOIN cor_themes_mob_img ctmi ON m.id_mobilier_image = ctmi.mob_img_id
         LEFT JOIN t_themes t ON t.id_theme = ctmi.theme_id
WHERE m.publication_status = 'DRAFT'
   OR m.publication_status = 'PENDING'
GROUP BY m.id_mobilier_image;

-- name: ValidatePendingMobilierImage :exec
UPDATE t_mobiliers_images
SET publication_status = 'PUBLISHED',
    publie             = true,
    parent_id          = NULL
WHERE id_mobilier_image = $1;

-- name: DeletePendingMobilierImage :exec
DELETE
FROM t_mobiliers_images
WHERE id_mobilier_image = $1;

-- name: CreateMobilierImage :one
INSERT INTO t_mobiliers_images
(titre_mob_img,
 description,
 historique,
 inscriptions,
 lieu_origine,
 lieu_conservation,
 bibliographie,
 protection,
 protection_commentaires,
 source,
 date_cr_ation,
 date_maj,
 contributeurs,
 id_commune,
 id_pays,
 publie,
 publication_status,
 parent_id)
VALUES (sqlc.arg(titre_mob_img),
        sqlc.arg(description),
        sqlc.arg(historique),
        sqlc.arg(inscriptions),
        sqlc.arg(origin),
        sqlc.arg(place),
        sqlc.arg(bibliographie),
        sqlc.arg(protection),
        sqlc.arg(protection_comment),
        sqlc.arg(source),
        NOW(),
        NOW(),
        sqlc.arg(contributors),
        sqlc.arg(id_commune),
        sqlc.arg(id_pays),
        false,
        'DRAFT',
        sqlc.arg(parent_id))
RETURNING id_mobilier_image;

-- name: AttachSieclesToMobImg :exec
INSERT INTO cor_siecles_mob_img
    (siecle_mob_img_id, mobilier_image_id)
SELECT unnest(sqlc.arg(siecle_id)::int[]), sqlc.arg(id);

-- name: AttachMediasToMobImg :exec
INSERT INTO cor_medias_mob_img
    (media_mob_img_id, mobilier_image_id)
SELECT unnest(sqlc.arg(media_ids)::int[]), sqlc.arg(id);

-- name: AttachThemesToMobImg :exec
INSERT INTO cor_themes_mob_img (theme_id, mob_img_id)
SELECT unnest(sqlc.arg(theme_ids)::int[]), sqlc.arg(id);

-- name: AttachNaturesToMobImg :exec
INSERT INTO cor_natures_mob_img (nature_id, mobilier_image_id)
SELECT unnest(sqlc.arg(nature_ids)::int[]), sqlc.arg(id);

-- name: AttachEtatsToMobImg :exec
INSERT INTO cor_etat_cons_mob_img
    (etat_cons_mob_img_id, mobilier_image_id)
SELECT unnest(sqlc.arg(etat_ids)::int[]), sqlc.arg(id);

-- name: AttachAuthorToMobImg :exec
INSERT INTO cor_auteur_fiche_mob_img
    (auteur_fiche_mob_img_id, mobilier_image_id)
VALUES (sqlc.arg(auteur_id), sqlc.arg(id));

-- name: AttachMateriauxToMobImg :exec
INSERT INTO cor_materiaux_mob_img
    (materiau_mob_img_id, mobilier_image_id)
SELECT unnest(sqlc.arg(materiau_ids)::int[]), sqlc.arg(id);

-- name: AttachTechniquesToMobImg :exec
INSERT INTO cor_techniques_mob_img
    (technique_id, mobilier_image_id)
SELECT unnest(sqlc.arg(techniques_ids)::int[]), sqlc.arg(id);

-- name: LinkMobImgToMonuLieu :exec
INSERT INTO cor_monu_lieu_mob_img
    (mobilier_image_id, monument_lieu_id)
SELECT sqlc.arg(id), unnest(sqlc.arg(monu_lieu_ids)::int[]);

-- name: LinkMobImgToPersMo :exec
INSERT INTO cor_mob_img_pers_mo
    (mobilier_image_id, pers_morale_id)
SELECT sqlc.arg(id), unnest(sqlc.arg(perso_mo_ids)::int[]);

-- name: LinkMobImgToPersPhy :exec
INSERT INTO cor_mob_img_pers_phy
    (mobilier_image_id, pers_physique_id)
SELECT sqlc.arg(id), unnest(sqlc.arg(perso_phy_ids)::int[]);

-- name: DetachSieclesFromMobImg :exec
DELETE
FROM cor_siecles_mob_img
WHERE mobilier_image_id = sqlc.arg(id);

-- name: DetachMediasFromMobImg :exec
DELETE
FROM cor_medias_mob_img
WHERE mobilier_image_id = sqlc.arg(id);

-- name: DetachThemesFromMobImg :exec
DELETE
FROM cor_themes_mob_img
WHERE mob_img_id = sqlc.arg(id);

-- name: DetachNaturesFromMobImg :exec
DELETE
FROM cor_natures_mob_img
WHERE mobilier_image_id = sqlc.arg(id);

-- name: DetachEtatsFromMobImg :exec
DELETE
FROM cor_etat_cons_mob_img
WHERE mobilier_image_id = sqlc.arg(id);

-- name: DetachAuthorFromMobImg :exec
DELETE
FROM cor_auteur_fiche_mob_img
WHERE mobilier_image_id = sqlc.arg(id);

-- name: DetachMateriauxFromMobImg :exec
DELETE
FROM cor_materiaux_mob_img
WHERE mobilier_image_id = sqlc.arg(id);

-- name: DetachTechniquesFromMobImg :exec
DELETE
FROM cor_techniques_mob_img
WHERE mobilier_image_id = sqlc.arg(id);

-- name: UnlinkMobImgFromMonuLieu :exec
DELETE
FROM cor_monu_lieu_mob_img
WHERE mobilier_image_id = sqlc.arg(id);

-- name: UnlinkMobImgFromPersMo :exec
DELETE
FROM cor_mob_img_pers_mo
WHERE mobilier_image_id = sqlc.arg(id);

-- name: UnlinkMobImgFromPersPhy :exec
DELETE
FROM cor_mob_img_pers_phy
WHERE mobilier_image_id = sqlc.arg(id);