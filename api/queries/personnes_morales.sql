-- name: GetFilteredPersonnesMorales :many
SELECT p.id_pers_morale AS id,
       p.titre_pers_mo  AS title,
       COALESCE(array_agg(DISTINCT bs.siecle_list) FILTER (WHERE bs.siecle_list IS NOT NULL),
                '{}')   AS siecles,
       COALESCE(array_agg(DISTINCT bpn.pers_mo_nature_type) FILTER (WHERE bpn.pers_mo_nature_type IS NOT NULL),
                '{}')   AS natures,
       COALESCE(array_agg(DISTINCT tm.chemin_media) FILTER (WHERE tm.chemin_media IS NOT NULL),
                '{}')   AS medias
FROM t_pers_morales p
         LEFT JOIN cor_siecles_pers_mo csp ON p.id_pers_morale = csp.pers_morale_id
         LEFT JOIN bib_siecle bs ON bs.id_siecle = csp.siecle_pers_mo_id

         LEFT JOIN cor_natures_pers_mo cnp ON p.id_pers_morale = cnp.pers_morale_id
         LEFT JOIN bib_pers_mo_natures bpn ON bpn.id_pers_mo_nature = cnp.pers_mo_nature_id

         LEFT JOIN cor_medias_mob_img cme ON p.id_pers_morale = cme.mobilier_image_id
         LEFT JOIN t_medias tm ON tm.id_media = cme.media_mob_img_id

         LEFT JOIN loc_communes c ON p.id_commune = c.id_commune
         LEFT JOIN loc_departements d ON c.id_departement = d.id_departement
         LEFT JOIN loc_regions r ON d.id_region = r.id_region
         LEFT JOIN loc_pays pa ON r.id_pays = pa.id_pays

         LEFT JOIN cor_auteur_fiche_pers_mo caf ON p.id_pers_morale = caf.pers_morale_id
         LEFT JOIN cor_mob_img_pers_mo cpm ON p.id_pers_morale = cpm.pers_morale_id
         LEFT JOIN cor_pers_phy_pers_mo cpp ON p.id_pers_morale = cpp.pers_morale_id
         LEFT JOIN cor_monu_lieu_pers_mo cmi ON p.id_pers_morale = cmi.pers_morale_id
WHERE (cardinality(sqlc.arg(siecles)::int[]) = 0 OR csp.siecle_pers_mo_id = ANY (sqlc.arg(siecles)::int[]))
   OR (cardinality(sqlc.arg(natures)::int[]) = 0 OR cnp.pers_mo_nature_id = ANY (sqlc.arg(natures)::int[]))
   OR (cardinality(sqlc.arg(auteurs_fiche)::int[]) = 0 OR
       caf.auteur_fiche_pers_mo_id = ANY (sqlc.arg(auteurs_fiche)::int[]))
   OR (cardinality(sqlc.arg(pers_phy)::int[]) = 0 OR cpp.pers_physique_id = ANY (sqlc.arg(pers_phy)::int[]))
   OR (cardinality(sqlc.arg(places)::int[]) = 0 OR cmi.monument_lieu_id = ANY (sqlc.arg(places)::int[]))
   OR (cardinality(sqlc.arg(furniture)::int[]) = 0 OR cpm.mobilier_image_id = ANY (sqlc.arg(furniture)::int[]))
   OR (cardinality(sqlc.arg(cities)::int[]) = 0 OR c.id_commune = ANY (sqlc.arg(cities)::int[]))
   OR (cardinality(sqlc.arg(departments)::int[]) = 0 OR d.id_departement = ANY (sqlc.arg(departments)::int[]))
   OR (cardinality(sqlc.arg(regions)::int[]) = 0 OR r.id_region = ANY (sqlc.arg(regions)::int[]))
   OR (cardinality(sqlc.arg(pays)::int[]) = 0 OR pa.id_pays = ANY (sqlc.arg(pays)::int[]))
GROUP BY p.id_pers_morale
ORDER BY p.id_pers_morale
LIMIT sqlc.arg(limit_param) OFFSET sqlc.arg(offset_param);

-- name: GetPersonneMoraleByID :one
SELECT p.id_pers_morale  AS id,
       p.titre_pers_mo   AS title,
       p.acte_fondation  AS foundation_deed,
       p.historique,
       p.bibliographie,
       p.simple_mention,
       p.fonctionnement,
       p.participation_vie_soc,
       p.objets,
       p.sources,
       p.date_creation,
       p.date_maj,
       p.publie,
       p.contributeurs,
       p.commentaires,
       p.user_id,
       p.biens,
       p.date_premiere_mention,
       p.date_derniere_mention,
       p.refondation_date,
       p.date_fin,
       p.texte_statuts,
       p.origine_sociale_prof,
       p.membres_connus,
       p.frequence_reunions,
       p.participation_vie_pol,
       p.funerailles,
       p.autres_fetes,
       p.inhumation_costume,
       p.fondateurs,
       p.statuts,
       p.autorisations,
       -- Redacteurs (auteurs fiche)
       COALESCE(
                       jsonb_agg(
                       DISTINCT jsonb_build_object(
                               'id', baf.id_auteur_fiche,
                               'name', baf.auteur_fiche_nom
                                )
                                ) FILTER (WHERE baf.id_auteur_fiche IS NOT NULL),
                       '[]'::jsonb
       ) AS authors,
       -- Commune
       jsonb_build_object(
               'id', c.id_commune,
               'name', c.nom_commune
       ) AS city,

       -- Département
       jsonb_build_object(
               'id', d.id_departement,
               'name', d.nom_departement
       ) AS department,

       -- Région
       jsonb_build_object(
               'id', r.id_region,
               'name', r.nom_region
       ) AS region,

       -- Pays
       jsonb_build_object(
               'id', pa.id_pays,
               'name', pa.nom_pays
       ) AS country,
       -- Natures
       COALESCE(
                       jsonb_agg(
                       DISTINCT jsonb_build_object(
                               'id', bpn.id_pers_mo_nature,
                               'name', bpn.pers_mo_nature_type
                                )
                                ) FILTER (WHERE bpn.id_pers_mo_nature IS NOT NULL),
                       '[]'::jsonb
       ) AS natures,
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
       )                 AS medias,
       -- Monuments lieux (IDs uniquement)
       COALESCE(array_agg(DISTINCT cml.monument_lieu_id) FILTER (WHERE cml.monument_lieu_id IS NOT NULL),
                '{}')    AS monuments_lieux_liees,
       -- Mobiliers images (IDs uniquement)
       COALESCE(array_agg(DISTINCT cpm.mobilier_image_id) FILTER (WHERE cpm.mobilier_image_id IS NOT NULL),
                '{}')    AS mobiliers_images_liees,
       -- Personnes physiques (IDs uniquement)
       COALESCE(array_agg(DISTINCT cpp.pers_physique_id) FILTER (WHERE cpp.pers_physique_id IS NOT NULL),
                '{}')    AS personnes_physiques_liees,
       -- Personnes morales liées (self-link)
       COALESCE(array_agg(DISTINCT
           CASE WHEN cpmpm.pers_morale_id_1 = p.id_pers_morale
                THEN cpmpm.pers_morale_id_2
                ELSE cpmpm.pers_morale_id_1
           END) FILTER (WHERE cpmpm.pers_morale_id_1 IS NOT NULL),
                '{}')    AS personnes_morales_liees,
       -- Siècles
       COALESCE(
                       jsonb_agg(
                       DISTINCT jsonb_build_object(
                               'id', bs.id_siecle,
                               'name', bs.siecle_list
                                )
                                ) FILTER (WHERE bs.id_siecle IS NOT NULL),
                       '[]'::jsonb
       ) AS centuries,
       -- Themes
       COALESCE(
                       jsonb_agg(
                       DISTINCT jsonb_build_object(
                               'id', t.id_theme,
                               'name', t.theme_type
                                )
                                ) FILTER (WHERE t.id_theme IS NOT NULL),
                       '[]'::jsonb
       ) AS themes,
       publication_status,
       parent_id
FROM t_pers_morales p
         LEFT JOIN cor_auteur_fiche_pers_mo cap ON p.id_pers_morale = cap.pers_morale_id
         LEFT JOIN bib_auteurs baf ON cap.auteur_fiche_pers_mo_id = baf.id_auteur_fiche
         LEFT JOIN loc_communes c ON p.id_commune = c.id_commune
         LEFT JOIN loc_departements d ON d.id_departement = COALESCE(p.id_departement, c.id_departement)
         LEFT JOIN loc_regions r ON r.id_region = COALESCE(p.id_region, d.id_region)
         LEFT JOIN loc_pays pa ON pa.id_pays = COALESCE(p.id_pays, r.id_pays)
         LEFT JOIN cor_natures_pers_mo cnp ON p.id_pers_morale = cnp.pers_morale_id
         LEFT JOIN bib_pers_mo_natures bpn ON cnp.pers_mo_nature_id = bpn.id_pers_mo_nature
         LEFT JOIN cor_medias_pers_mo cmp ON p.id_pers_morale = cmp.pers_morale_id
         LEFT JOIN t_medias tm ON cmp.media_pers_mo_id = tm.id_media
         LEFT JOIN cor_monu_lieu_pers_mo cml ON p.id_pers_morale = cml.pers_morale_id
         LEFT JOIN cor_mob_img_pers_mo cpm ON p.id_pers_morale = cpm.pers_morale_id
         LEFT JOIN cor_pers_phy_pers_mo cpp ON p.id_pers_morale = cpp.pers_morale_id
         LEFT JOIN cor_pers_mo_pers_mo cpmpm ON p.id_pers_morale = cpmpm.pers_morale_id_1 OR p.id_pers_morale = cpmpm.pers_morale_id_2
         LEFT JOIN cor_siecles_pers_mo csp ON p.id_pers_morale = csp.pers_morale_id
         LEFT JOIN bib_siecle bs ON csp.siecle_pers_mo_id = bs.id_siecle
         LEFT JOIN cor_themes_pers_mo ctpm ON p.id_pers_morale = ctpm.pers_mo_id
         LEFT JOIN t_themes t ON t.id_theme = ctpm.theme_id
WHERE p.id_pers_morale = $1
GROUP BY p.id_pers_morale,
         c.id_commune, c.nom_commune,
         d.id_departement, d.nom_departement,
         r.id_region, r.nom_region,
         pa.id_pays, pa.nom_pays;

-- name: GetPendingPersonnesMorales :many
SELECT p.id_pers_morale       AS id,
       p.titre_pers_mo        AS title,
       p.acte_fondation       AS foundation_deed,
       p.historique,
       p.bibliographie,
       p.simple_mention,
       p.fonctionnement,
       p.participation_vie_soc,
       p.objets,
       p.sources,
       p.date_creation,
       p.date_maj,
       p.publie,
       p.contributeurs,
       p.commentaires,
       p.parent_id,
       p.biens,
       p.date_premiere_mention,
       p.date_derniere_mention,
       p.refondation_date,
       p.date_fin,
       p.texte_statuts,
       p.origine_sociale_prof,
       p.membres_connus,
       p.frequence_reunions,
       p.participation_vie_pol,
       p.funerailles,
       p.autres_fetes,
       p.inhumation_costume,
       p.fondateurs,
       p.statuts,
       p.autorisations,
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
       -- Natures
       COALESCE(array_agg(DISTINCT bpn.pers_mo_nature_type) FILTER (WHERE bpn.pers_mo_nature_type IS NOT NULL),
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
       COALESCE(array_agg(DISTINCT cml.monument_lieu_id) FILTER (WHERE cml.monument_lieu_id IS NOT NULL),
                '{}')         AS monuments_lieux_liees,
       -- Mobiliers images (IDs uniquement)
       COALESCE(array_agg(DISTINCT cpm.mobilier_image_id) FILTER (WHERE cpm.mobilier_image_id IS NOT NULL),
                '{}')         AS mobiliers_images_liees,
       -- Personnes physiques (IDs uniquement)
       COALESCE(array_agg(DISTINCT cpp.pers_physique_id) FILTER (WHERE cpp.pers_physique_id IS NOT NULL),
                '{}')         AS personnes_physiques_liees,
       -- Personnes morales liées (self-link)
       COALESCE(array_agg(DISTINCT
           CASE WHEN cpmpm.pers_morale_id_1 = p.id_pers_morale
                THEN cpmpm.pers_morale_id_2
                ELSE cpmpm.pers_morale_id_1
           END) FILTER (WHERE cpmpm.pers_morale_id_1 IS NOT NULL),
                '{}')         AS personnes_morales_liees,
       -- Siècles
       COALESCE(array_agg(DISTINCT bs.siecle_list) FILTER (WHERE bs.siecle_list IS NOT NULL),
                '{}')         AS siecles,
       -- Themes
       COALESCE(array_agg(DISTINCT t.theme_type) FILTER (WHERE t.theme_type IS NOT NULL),
                '{}')         AS themes
FROM t_pers_morales p
         LEFT JOIN cor_auteur_fiche_pers_mo cap ON p.id_pers_morale = cap.pers_morale_id
         LEFT JOIN bib_auteurs baf ON cap.auteur_fiche_pers_mo_id = baf.id_auteur_fiche
         LEFT JOIN loc_communes c ON p.id_commune = c.id_commune
         LEFT JOIN loc_departements d ON c.id_departement = d.id_departement
         LEFT JOIN loc_regions r ON d.id_region = r.id_region
         LEFT JOIN loc_pays pa ON r.id_pays = p.id_pays
         LEFT JOIN cor_medias_pers_mo cmp ON p.id_pers_morale = cmp.pers_morale_id
         LEFT JOIN cor_natures_pers_mo cnp ON p.id_pers_morale = cnp.pers_morale_id
         LEFT JOIN bib_pers_mo_natures bpn ON cnp.pers_mo_nature_id = bpn.id_pers_mo_nature
         LEFT JOIN t_medias tm ON cmp.media_pers_mo_id = tm.id_media
         LEFT JOIN cor_monu_lieu_pers_mo cml ON p.id_pers_morale = cml.pers_morale_id
         LEFT JOIN cor_mob_img_pers_mo cpm ON p.id_pers_morale = cpm.pers_morale_id
         LEFT JOIN cor_pers_phy_pers_mo cpp ON p.id_pers_morale = cpp.pers_morale_id
         LEFT JOIN cor_pers_mo_pers_mo cpmpm ON p.id_pers_morale = cpmpm.pers_morale_id_1 OR p.id_pers_morale = cpmpm.pers_morale_id_2
         LEFT JOIN cor_siecles_pers_mo csp ON p.id_pers_morale = csp.pers_morale_id
         LEFT JOIN bib_siecle bs ON csp.siecle_pers_mo_id = bs.id_siecle
         LEFT JOIN cor_themes_pers_mo ctpm ON p.id_pers_morale = ctpm.pers_mo_id
         LEFT JOIN t_themes t ON t.id_theme = ctpm.theme_id
WHERE p.publication_status = 'PENDING'
GROUP BY p.id_pers_morale;

-- name: GetDraftPersonnesMorales :many
SELECT p.id_pers_morale       AS id,
       p.titre_pers_mo        AS title,
       p.acte_fondation       AS foundation_deed,
       p.historique,
       p.bibliographie,
       p.simple_mention,
       p.fonctionnement,
       p.participation_vie_soc,
       p.objets,
       p.sources,
       p.date_creation,
       p.date_maj,
       p.publie,
       p.contributeurs,
       p.commentaires,
       p.parent_id,
       p.biens,
       p.date_premiere_mention,
       p.date_derniere_mention,
       p.refondation_date,
       p.date_fin,
       p.texte_statuts,
       p.origine_sociale_prof,
       p.membres_connus,
       p.frequence_reunions,
       p.participation_vie_pol,
       p.funerailles,
       p.autres_fetes,
       p.inhumation_costume,
       p.fondateurs,
       p.statuts,
       p.autorisations,
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
       -- Natures
       COALESCE(array_agg(DISTINCT bpn.pers_mo_nature_type) FILTER (WHERE bpn.pers_mo_nature_type IS NOT NULL),
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
       COALESCE(array_agg(DISTINCT cml.monument_lieu_id) FILTER (WHERE cml.monument_lieu_id IS NOT NULL),
                '{}')         AS monuments_lieux_liees,
       -- Mobiliers images (IDs uniquement)
       COALESCE(array_agg(DISTINCT cpm.mobilier_image_id) FILTER (WHERE cpm.mobilier_image_id IS NOT NULL),
                '{}')         AS mobiliers_images_liees,
       -- Personnes physiques (IDs uniquement)
       COALESCE(array_agg(DISTINCT cpp.pers_physique_id) FILTER (WHERE cpp.pers_physique_id IS NOT NULL),
                '{}')         AS personnes_physiques_liees,
       -- Personnes morales liées (self-link)
       COALESCE(array_agg(DISTINCT
           CASE WHEN cpmpm.pers_morale_id_1 = p.id_pers_morale
                THEN cpmpm.pers_morale_id_2
                ELSE cpmpm.pers_morale_id_1
           END) FILTER (WHERE cpmpm.pers_morale_id_1 IS NOT NULL),
                '{}')         AS personnes_morales_liees,
       -- Siècles
       COALESCE(array_agg(DISTINCT bs.siecle_list) FILTER (WHERE bs.siecle_list IS NOT NULL),
                '{}')         AS siecles,
       -- Themes
       COALESCE(array_agg(DISTINCT t.theme_type) FILTER (WHERE t.theme_type IS NOT NULL),
                '{}')         AS themes
FROM t_pers_morales p
         LEFT JOIN cor_auteur_fiche_pers_mo cap ON p.id_pers_morale = cap.pers_morale_id
         LEFT JOIN bib_auteurs baf ON cap.auteur_fiche_pers_mo_id = baf.id_auteur_fiche
         LEFT JOIN loc_communes c ON p.id_commune = c.id_commune
         LEFT JOIN loc_departements d ON c.id_departement = d.id_departement
         LEFT JOIN loc_regions r ON d.id_region = r.id_region
         LEFT JOIN loc_pays pa ON r.id_pays = p.id_pays
         LEFT JOIN cor_medias_pers_mo cmp ON p.id_pers_morale = cmp.pers_morale_id
         LEFT JOIN cor_natures_pers_mo cnp ON p.id_pers_morale = cnp.pers_morale_id
         LEFT JOIN bib_pers_mo_natures bpn ON cnp.pers_mo_nature_id = bpn.id_pers_mo_nature
         LEFT JOIN t_medias tm ON cmp.media_pers_mo_id = tm.id_media
         LEFT JOIN cor_monu_lieu_pers_mo cml ON p.id_pers_morale = cml.pers_morale_id
         LEFT JOIN cor_mob_img_pers_mo cpm ON p.id_pers_morale = cpm.pers_morale_id
         LEFT JOIN cor_pers_phy_pers_mo cpp ON p.id_pers_morale = cpp.pers_morale_id
         LEFT JOIN cor_pers_mo_pers_mo cpmpm ON p.id_pers_morale = cpmpm.pers_morale_id_1 OR p.id_pers_morale = cpmpm.pers_morale_id_2
         LEFT JOIN cor_siecles_pers_mo csp ON p.id_pers_morale = csp.pers_morale_id
         LEFT JOIN bib_siecle bs ON csp.siecle_pers_mo_id = bs.id_siecle
         LEFT JOIN cor_themes_pers_mo ctpm ON p.id_pers_morale = ctpm.pers_mo_id
         LEFT JOIN t_themes t ON t.id_theme = ctpm.theme_id
WHERE p.publication_status = 'DRAFT'
and p.user_id = $1
GROUP BY p.id_pers_morale;

-- name: ValidatePendingPersonneMorales :exec
UPDATE t_pers_morales
SET publication_status = 'PUBLISHED',
    publie             = true,
    parent_id          = NULL
WHERE id_pers_morale = $1
AND publication_status = 'PENDING';

-- name: DeletePersonneMorale :exec
DELETE
FROM t_pers_morales
WHERE id_pers_morale = $1;

-- name: SubmitDraftPersonneMorale :exec
UPDATE t_pers_morales
SET publication_status = 'PENDING'
WHERE id_pers_morale = $1
  AND publication_status = 'DRAFT';

-- name: CreatePersMorale :one
INSERT INTO t_pers_morales
(titre_pers_mo,
 commentaires,
 historique,
 acte_fondation,
 simple_mention,
 texte_statuts,
 fonctionnement,
 participation_vie_soc,
 objets,
 bibliographie,
 sources,
 date_creation,
 date_maj,
 contributeurs,
 id_commune,
 id_departement,
 id_region,
 id_pays,
 publie,
 publication_status,
 parent_id,
 user_id,
 biens,
 date_premiere_mention,
 date_derniere_mention,
 refondation_date,
 date_fin,
 origine_sociale_prof,
 membres_connus,
 frequence_reunions,
 participation_vie_pol,
 funerailles,
 autres_fetes,
 inhumation_costume,
 fondateurs,
 statuts,
 autorisations)
VALUES (sqlc.arg(title),
        sqlc.arg(comment),
        sqlc.arg(historique),
        sqlc.arg(acte_fondation),
        sqlc.arg(simple_mention),
        sqlc.arg(texte_statuts),
        sqlc.arg(fonctionnement),
        sqlc.arg(participation_vie_soc),
        sqlc.arg(objets),
        sqlc.arg(bibliographie),
        sqlc.arg(sources),
        COALESCE(sqlc.arg(date_creation), NOW()),
        NOW(),
        sqlc.arg(contributeurs),
        sqlc.arg(id_commune),
        sqlc.arg(id_departement),
        sqlc.arg(id_region),
        sqlc.arg(id_pays),
        false,
        sqlc.arg(publication_status),
        sqlc.arg(parent_id),
        sqlc.arg(user_id),
        sqlc.arg(biens),
        sqlc.arg(date_premiere_mention),
        sqlc.arg(date_derniere_mention),
        sqlc.arg(refondation_date),
        sqlc.arg(date_fin),
        sqlc.arg(origine_sociale_prof),
        sqlc.arg(membres_connus),
        sqlc.arg(frequence_reunions),
        sqlc.arg(participation_vie_pol),
        sqlc.arg(funerailles),
        sqlc.arg(autres_fetes),
        sqlc.arg(inhumation_costume),
        sqlc.arg(fondateurs),
        sqlc.arg(statuts),
        sqlc.arg(autorisations))
RETURNING id_pers_morale;

-- name: AttachSieclesToPersMo :exec
INSERT INTO cor_siecles_pers_mo
    (siecle_pers_mo_id, pers_morale_id)
SELECT unnest(sqlc.arg(siecle_id)::int[]), sqlc.arg(id);

-- name: AttachMediasToPersMo :exec
INSERT INTO cor_medias_pers_mo
    (media_pers_mo_id, pers_morale_id)
SELECT unnest(sqlc.arg(media_ids)::int[]), sqlc.arg(id);

-- name: AttachThemesToPersMo :exec
INSERT INTO cor_themes_pers_mo (theme_id, pers_mo_id)
SELECT unnest(sqlc.arg(theme_ids)::int[]), sqlc.arg(id);

-- name: AttachNaturesToPersMo :exec
INSERT INTO cor_natures_pers_mo (pers_mo_nature_id, pers_morale_id)
SELECT unnest(sqlc.arg(nature_ids)::int[]), sqlc.arg(id);

-- name: AttachAuthorToPersMo :exec
INSERT INTO cor_auteur_fiche_pers_mo
    (auteur_fiche_pers_mo_id, pers_morale_id)
VALUES (sqlc.arg(auteur_id), sqlc.arg(id));

-- name: LinkPersMoToMonuLieu :exec
INSERT INTO cor_monu_lieu_pers_mo
    (pers_morale_id, monument_lieu_id)
SELECT sqlc.arg(id), unnest(sqlc.arg(monu_lieu_ids)::int[]);

-- name: LinkPersMoToMobImg :exec
INSERT INTO cor_mob_img_pers_mo
    (pers_morale_id, mobilier_image_id)
SELECT sqlc.arg(id), unnest(sqlc.arg(mob_img_ids)::int[]);

-- name: LinkPersMoToPersPhy :exec
INSERT INTO cor_pers_phy_pers_mo
    (pers_morale_id, pers_physique_id)
SELECT sqlc.arg(id), unnest(sqlc.arg(perso_phy_ids)::int[]);

-- name: DetachSieclesFromPersMo :exec
DELETE
FROM cor_siecles_pers_mo
WHERE pers_morale_id = sqlc.arg(id);

-- name: DetachMediasFromPersMo :exec
DELETE
FROM cor_medias_pers_mo
WHERE pers_morale_id = sqlc.arg(id);

-- name: DetachThemesFromPersMo :exec
DELETE
FROM cor_themes_pers_mo
WHERE pers_mo_id = sqlc.arg(id);

-- name: DetachNaturesFromPersMo :exec
DELETE
FROM cor_natures_pers_mo
WHERE pers_morale_id = sqlc.arg(id);

-- name: DetachAuthorFromPersMo :exec
DELETE
FROM cor_auteur_fiche_pers_mo
WHERE pers_morale_id = sqlc.arg(id);

-- name: UnlinkPersMoFromMonuLieu :exec
DELETE
FROM cor_monu_lieu_pers_mo
WHERE pers_morale_id = sqlc.arg(id);

-- name: UnlinkPersMoFromMobImg :exec
DELETE
FROM cor_mob_img_pers_mo
WHERE pers_morale_id = sqlc.arg(id);

-- name: UnlinkPersMoFromPersPhy :exec
DELETE
FROM cor_pers_phy_pers_mo
WHERE pers_morale_id = sqlc.arg(id);

-- name: LinkPersMoToPersMo :exec
INSERT INTO cor_pers_mo_pers_mo
    (pers_morale_id_1, pers_morale_id_2)
SELECT sqlc.arg(id), unnest(sqlc.arg(pers_mo_ids)::int[]);

-- name: UnlinkPersMoFromPersMo :exec
DELETE
FROM cor_pers_mo_pers_mo
WHERE pers_morale_id_1 = sqlc.arg(id)
   OR pers_morale_id_2 = sqlc.arg(id);