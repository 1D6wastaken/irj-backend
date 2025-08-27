-- name: SearchGlobal :many
WITH query AS (SELECT unnest(string_to_array(lower(sqlc.arg('q')), ' ')) AS term)
SELECT *, COUNT(*) OVER () AS total_count
FROM (
         -- ===============================
         -- 1. Monuments & lieux
         -- ===============================
         SELECT m.id_monument_lieu                           AS id,
                m.titre_monu_lieu                            AS title,
                COALESCE(array_agg(DISTINCT bs.siecle_list) FILTER (WHERE bs.siecle_list IS NOT NULL),
                         '{}')                               AS siecles,
                COALESCE(array_agg(DISTINCT bmn.monu_lieu_nature_type)
                         FILTER (WHERE bmn.monu_lieu_nature_type IS NOT NULL),
                         '{}')                               AS natures,
                COALESCE(
                                jsonb_agg(
                                DISTINCT jsonb_build_object(
                                        'id', tm.id_media,
                                        'titre', tm.titre_media
                                         )
                                         ) FILTER (
                                    WHERE tm.titre_media IS NOT NULL
                                AND tm.titre_media <> ''
                                AND tm.chemin_media IS NOT NULL
                                AND tm.chemin_media <> ''
                                AND jsonb_typeof(tm.chemin_media::jsonb) = 'array'
                                AND EXISTS (SELECT 1
                                            FROM jsonb_array_elements(tm.chemin_media::jsonb) AS elem
                                            WHERE COALESCE(elem ->> 'path', '') <> '')
                                    ),
                                '[]'::jsonb
                )                                            AS medias,
                '{}'::text[]                                 AS professions,
                'monuments_lieux'                            AS source,
                similarity(m.titre_monu_lieu, sqlc.arg('q')) AS score
         FROM t_monuments_lieux m
                  LEFT JOIN loc_communes c ON c.id_commune = m.id_commune
                  LEFT JOIN loc_departements d ON d.id_departement = c.id_departement
                  LEFT JOIN loc_regions r ON r.id_region = d.id_region
                  LEFT JOIN loc_pays p ON p.id_pays = r.id_pays
                  LEFT JOIN cor_siecles_monu_lieu csm ON csm.monument_lieu_id = m.id_monument_lieu
                  LEFT JOIN bib_siecle bs ON bs.id_siecle = csm.siecle_monu_lieu_id
                  LEFT JOIN cor_natures_monu_lieu cnm ON cnm.monument_lieu_id = m.id_monument_lieu
                  LEFT JOIN bib_monu_lieu_natures bmn ON bmn.id_monu_lieu_nature = cnm.monu_lieu_nature_id
                  LEFT JOIN cor_etat_cons_monu_lieu cem ON cem.monument_lieu_id = m.id_monument_lieu
                  LEFT JOIN cor_materiaux_monu_lieu cmm ON cmm.monument_lieu_id = m.id_monument_lieu
                  LEFT JOIN cor_medias_monu_lieu cme ON m.id_monument_lieu = cme.monument_lieu_id
                  LEFT JOIN t_medias tm ON tm.id_media = cme.media_monu_lieu_id
         WHERE sqlc.arg('include_monuments_lieux') = true
           AND (sqlc.arg('q') IS NULL OR m.titre_monu_lieu ILIKE '%' || sqlc.arg('q') || '%')
           AND ((sqlc.arg('siecles')::int[]) IS NULL OR cardinality(sqlc.arg('siecles')::int[]) = 0 OR
                csm.siecle_monu_lieu_id = ANY (sqlc.arg('siecles')::int[]))
           AND ((sqlc.arg('pays')::int[]) IS NULL OR cardinality(sqlc.arg('pays')::int[]) = 0 OR
                p.id_pays = ANY (sqlc.arg('pays')::int[]))
           AND ((sqlc.arg('region')::int[]) IS NULL OR cardinality(sqlc.arg('region')::int[]) = 0 OR
                r.id_region = ANY (sqlc.arg('region')::int[]))
           AND ((sqlc.arg('departement')::int[]) IS NULL OR cardinality(sqlc.arg('departement')::int[]) = 0 OR
                d.id_departement = ANY (sqlc.arg('departement')::int[]))
           AND ((sqlc.arg('commune')::int[]) IS NULL OR cardinality(sqlc.arg('commune')::int[]) = 0 OR
                c.id_commune = ANY (sqlc.arg('commune')::int[]))
           AND ((sqlc.arg('natures_monu')::int[]) IS NULL OR cardinality(sqlc.arg('natures_monu')::int[]) = 0 OR
                cnm.monu_lieu_nature_id = ANY (sqlc.arg('natures_monu')::int[]))
           AND ((sqlc.arg('etats_monu')::int[]) IS NULL OR cardinality(sqlc.arg('etats_monu')::int[]) = 0 OR
                cnm.monu_lieu_nature_id = ANY (sqlc.arg('etats_monu')::int[]))
           AND ((sqlc.arg('materiaux_monu')::int[]) IS NULL OR cardinality(sqlc.arg('materiaux_monu')::int[]) = 0 OR
                cnm.monu_lieu_nature_id = ANY (sqlc.arg('materiaux_monu')::int[]))
           AND m.publie = true AND m.publication_status = 'PUBLISHED'
         GROUP BY m.id_monument_lieu

         UNION ALL

         -- ===============================
         -- 2. Mobiliers & images
         -- ===============================
         SELECT mob.id_mobilier_image                        AS id,
                mob.titre_mob_img                            AS title,
                COALESCE(array_agg(DISTINCT bs.siecle_list) FILTER (WHERE bs.siecle_list IS NOT NULL),
                         '{}')                               AS siecles,
                COALESCE(array_agg(DISTINCT bmn.nature_type) FILTER (WHERE bmn.nature_type IS NOT NULL),
                         '{}')                               AS natures,
                COALESCE(
                                jsonb_agg(
                                DISTINCT jsonb_build_object(
                                        'id', tm.id_media,
                                        'titre', tm.titre_media
                                         )
                                         ) FILTER (
                                    WHERE tm.titre_media IS NOT NULL
                                AND tm.titre_media <> ''
                                AND tm.chemin_media IS NOT NULL
                                AND tm.chemin_media <> ''
                                AND jsonb_typeof(tm.chemin_media::jsonb) = 'array'
                                AND EXISTS (SELECT 1
                                            FROM jsonb_array_elements(tm.chemin_media::jsonb) AS elem
                                            WHERE COALESCE(elem ->> 'path', '') <> '')
                                    ),
                                '[]'::jsonb
                )                                            AS medias,
                '{}'::text[]                                 AS professions,
                'mobiliers_images'                           AS source,
                similarity(mob.titre_mob_img, sqlc.arg('q')) AS score
         FROM t_mobiliers_images mob
                  LEFT JOIN loc_communes c ON c.id_commune = mob.id_commune
                  LEFT JOIN loc_departements d ON d.id_departement = c.id_departement
                  LEFT JOIN loc_regions r ON r.id_region = d.id_region
                  LEFT JOIN loc_pays p ON p.id_pays = r.id_pays
                  LEFT JOIN cor_siecles_mob_img csm ON csm.mobilier_image_id = mob.id_mobilier_image
                  LEFT JOIN bib_siecle bs ON bs.id_siecle = csm.siecle_mob_img_id
                  LEFT JOIN cor_natures_mob_img cnm ON cnm.mobilier_image_id = mob.id_mobilier_image
                  LEFT JOIN bib_mob_img_natures bmn ON bmn.id_nature = cnm.nature_id
                  LEFT JOIN cor_etat_cons_mob_img cem ON cem.mobilier_image_id = mob.id_mobilier_image
                  LEFT JOIN cor_materiaux_mob_img cmm ON cmm.mobilier_image_id = mob.id_mobilier_image
                  LEFT JOIN cor_techniques_mob_img ctm ON ctm.mobilier_image_id = mob.id_mobilier_image
                  LEFT JOIN cor_medias_mob_img cme ON mob.id_mobilier_image = cme.mobilier_image_id
                  LEFT JOIN t_medias tm ON tm.id_media = cme.media_mob_img_id
         WHERE sqlc.arg('include_mobiliers_images') = true
           AND (sqlc.arg('q') IS NULL OR mob.titre_mob_img ILIKE '%' || sqlc.arg('q') || '%')
           AND ((sqlc.arg('siecles')::int[]) IS NULL OR cardinality(sqlc.arg('siecles')::int[]) = 0 OR
                csm.siecle_mob_img_id = ANY (sqlc.arg('siecles')::int[]))
           AND ((sqlc.arg('pays')::int[]) IS NULL OR cardinality(sqlc.arg('pays')::int[]) = 0 OR
                p.id_pays = ANY (sqlc.arg('pays')::int[]))
           AND ((sqlc.arg('region')::int[]) IS NULL OR cardinality(sqlc.arg('region')::int[]) = 0 OR
                r.id_region = ANY (sqlc.arg('region')::int[]))
           AND ((sqlc.arg('departement')::int[]) IS NULL OR cardinality(sqlc.arg('departement')::int[]) = 0 OR
                d.id_departement = ANY (sqlc.arg('departement')::int[]))
           AND ((sqlc.arg('commune')::int[]) IS NULL OR cardinality(sqlc.arg('commune')::int[]) = 0 OR
                c.id_commune = ANY (sqlc.arg('commune')::int[]))
           AND ((sqlc.arg('natures_mob')::int[]) IS NULL OR cardinality(sqlc.arg('natures_mob')::int[]) = 0 OR
                cnm.nature_id = ANY (sqlc.arg('natures_mob')::int[]))
           AND ((sqlc.arg('etats_mob')::int[]) IS NULL OR cardinality(sqlc.arg('etats_mob')::int[]) = 0 OR
                cem.etat_cons_mob_img_id = ANY (sqlc.arg('etats_mob')::int[]))
           AND ((sqlc.arg('materiaux_mob')::int[]) IS NULL OR cardinality(sqlc.arg('materiaux_mob')::int[]) = 0 OR
                cmm.materiau_mob_img_id = ANY (sqlc.arg('materiaux_mob')::int[]))
           AND ((sqlc.arg('techniques_mob')::int[]) IS NULL OR cardinality(sqlc.arg('techniques_mob')::int[]) = 0 OR
                ctm.technique_id = ANY (sqlc.arg('techniques_mob')::int[]))
           AND mob.publie = true AND mob.publication_status = 'PUBLISHED'
         GROUP BY mob.id_mobilier_image

         UNION ALL

         -- ===============================
         -- 3. Personnes morales
         -- ===============================
         SELECT pm.id_pers_morale                           AS id,
                pm.titre_pers_mo                            AS title,
                COALESCE(array_agg(DISTINCT bs.siecle_list) FILTER (WHERE bs.siecle_list IS NOT NULL),
                         '{}')                              AS siecles,
                COALESCE(array_agg(DISTINCT bpn.pers_mo_nature_type) FILTER (WHERE bpn.pers_mo_nature_type IS NOT NULL),
                         '{}')                              AS natures,
                COALESCE(
                                jsonb_agg(
                                DISTINCT jsonb_build_object(
                                        'id', tm.id_media,
                                        'titre', tm.titre_media
                                         )
                                         ) FILTER (
                                    WHERE tm.titre_media IS NOT NULL
                                AND tm.titre_media <> ''
                                AND tm.chemin_media IS NOT NULL
                                AND tm.chemin_media <> ''
                                AND jsonb_typeof(tm.chemin_media::jsonb) = 'array'
                                AND EXISTS (SELECT 1
                                            FROM jsonb_array_elements(tm.chemin_media::jsonb) AS elem
                                            WHERE COALESCE(elem ->> 'path', '') <> '')
                                    ),
                                '[]'::jsonb
                )                                           AS medias,
                '{}'::text[]                                AS professions,
                'personnes_morales'                         AS source,
                similarity(pm.titre_pers_mo, sqlc.arg('q')) AS score
         FROM t_pers_morales pm
                  LEFT JOIN loc_communes c ON c.id_commune = pm.id_commune
                  LEFT JOIN loc_departements d ON d.id_departement = c.id_departement
                  LEFT JOIN loc_regions r ON r.id_region = d.id_region
                  LEFT JOIN loc_pays p ON p.id_pays = r.id_pays
                  LEFT JOIN cor_siecles_pers_mo csp ON csp.pers_morale_id = pm.id_pers_morale
                  LEFT JOIN bib_siecle bs ON bs.id_siecle = csp.siecle_pers_mo_id
                  LEFT JOIN cor_natures_pers_mo cnp ON cnp.pers_morale_id = pm.id_pers_morale
                  LEFT JOIN bib_pers_mo_natures bpn ON bpn.id_pers_mo_nature = cnp.pers_mo_nature_id
                  LEFT JOIN cor_medias_mob_img cme ON pm.id_pers_morale = cme.mobilier_image_id
                  LEFT JOIN t_medias tm ON tm.id_media = cme.media_mob_img_id
         WHERE sqlc.arg('include_pers_morales') = true
           AND (sqlc.arg('q') IS NULL OR pm.titre_pers_mo ILIKE '%' || sqlc.arg('q') || '%')
           AND ((sqlc.arg('siecles')::int[]) IS NULL OR cardinality(sqlc.arg('siecles')::int[]) = 0 OR
                csp.siecle_pers_mo_id = ANY (sqlc.arg('siecles')::int[]))
           AND ((sqlc.arg('pays')::int[]) IS NULL OR cardinality(sqlc.arg('pays')::int[]) = 0 OR
                p.id_pays = ANY (sqlc.arg('pays')::int[]))
           AND ((sqlc.arg('region')::int[]) IS NULL OR cardinality(sqlc.arg('region')::int[]) = 0 OR
                r.id_region = ANY (sqlc.arg('region')::int[]))
           AND ((sqlc.arg('departement')::int[]) IS NULL OR cardinality(sqlc.arg('departement')::int[]) = 0 OR
                d.id_departement = ANY (sqlc.arg('departement')::int[]))
           AND ((sqlc.arg('commune')::int[]) IS NULL OR cardinality(sqlc.arg('commune')::int[]) = 0 OR
                c.id_commune = ANY (sqlc.arg('commune')::int[]))
           AND ((sqlc.arg('natures_pers_mo')::int[]) IS NULL OR cardinality(sqlc.arg('natures_pers_mo')::int[]) = 0 OR
                cnp.pers_mo_nature_id = ANY (sqlc.arg('natures_pers_mo')::int[]))
           AND pm.publie = true AND pm.publication_status = 'PUBLISHED'
         GROUP BY pm.id_pers_morale

         UNION ALL

         -- ===============================
         -- 4. Personnes physiques
         -- ===============================
         SELECT pp.id_pers_physique                               AS id,
                pp.prenom_nom_pers_phy                            AS title,
                COALESCE(array_agg(DISTINCT bs.siecle_list) FILTER (WHERE bs.siecle_list IS NOT NULL),
                         '{}')                                    AS siecles,
                '{}'::text[]                                      AS natures,
                COALESCE(
                                jsonb_agg(
                                DISTINCT jsonb_build_object(
                                        'id', tm.id_media,
                                        'titre', tm.titre_media
                                         )
                                         ) FILTER (
                                    WHERE tm.titre_media IS NOT NULL
                                AND tm.titre_media <> ''
                                AND tm.chemin_media IS NOT NULL
                                AND tm.chemin_media <> ''
                                AND jsonb_typeof(tm.chemin_media::jsonb) = 'array'
                                AND EXISTS (SELECT 1
                                            FROM jsonb_array_elements(tm.chemin_media::jsonb) AS elem
                                            WHERE COALESCE(elem ->> 'path', '') <> '')
                                    ),
                                '[]'::jsonb
                )                                                 AS medias,
                COALESCE(array_agg(DISTINCT bpp.profession_type) FILTER (WHERE bpp.profession_type IS NOT NULL),
                         '{}')                                    AS professions,
                'personnes_physiques'                             AS source,
                similarity(pp.prenom_nom_pers_phy, sqlc.arg('q')) AS score
         FROM t_pers_physiques pp
                  LEFT JOIN loc_communes c ON c.id_commune = pp.id_commune
                  LEFT JOIN loc_departements d ON d.id_departement = c.id_departement
                  LEFT JOIN loc_regions r ON r.id_region = d.id_region
                  LEFT JOIN loc_pays p ON p.id_pays = r.id_pays
                  LEFT JOIN cor_siecles_pers_phy csp ON csp.pers_physique_id = pp.id_pers_physique
                  LEFT JOIN bib_siecle bs ON bs.id_siecle = csp.siecle_pers_phy_id
                  LEFT JOIN cor_professions_pers_phy cpp ON cpp.pers_physique_id = pp.id_pers_physique
                  LEFT JOIN bib_pers_phy_professions bpp ON bpp.id_profession = cpp.profession_id
                  LEFT JOIN cor_modes_deplacements_pers_phy cmd ON cmd.pers_physique_id = pp.id_pers_physique
                  LEFT JOIN cor_medias_pers_phy cmp ON pp.id_pers_physique = cmp.pers_physique_id
                  LEFT JOIN t_medias tm ON tm.id_media = cmp.media_pers_phy_id
         WHERE sqlc.arg('include_pers_physiques') = true
           AND (sqlc.arg('q') IS NULL OR pp.prenom_nom_pers_phy ILIKE '%' || sqlc.arg('q') || '%')
           AND ((sqlc.arg('siecles')::int[]) IS NULL OR cardinality(sqlc.arg('siecles')::int[]) = 0 OR
                csp.siecle_pers_phy_id = ANY (sqlc.arg('siecles')::int[]))
           AND ((sqlc.arg('pays')::int[]) IS NULL OR cardinality(sqlc.arg('pays')::int[]) = 0 OR
                p.id_pays = ANY (sqlc.arg('pays')::int[]))
           AND ((sqlc.arg('region')::int[]) IS NULL OR cardinality(sqlc.arg('region')::int[]) = 0 OR
                r.id_region = ANY (sqlc.arg('region')::int[]))
           AND ((sqlc.arg('departement')::int[]) IS NULL OR cardinality(sqlc.arg('departement')::int[]) = 0 OR
                d.id_departement = ANY (sqlc.arg('departement')::int[]))
           AND ((sqlc.arg('commune')::int[]) IS NULL OR cardinality(sqlc.arg('commune')::int[]) = 0 OR
                c.id_commune = ANY (sqlc.arg('commune')::int[]))
           AND ((sqlc.arg('professions')::int[]) IS NULL OR cardinality(sqlc.arg('professions')::int[]) = 0 OR
                cpp.profession_id = ANY (sqlc.arg('professions')::int[]))
           AND ((sqlc.arg('modes_deplacements')::int[]) IS NULL OR
                cardinality(sqlc.arg('modes_deplacements')::int[]) = 0 OR
                cmd.mode_deplacement_id = ANY (sqlc.arg('modes_deplacements')::int[]))
           AND pp.publie = true AND pp.publication_status = 'PUBLISHED'
         GROUP BY pp.id_pers_physique) AS results


ORDER BY score DESC
LIMIT sqlc.arg(limit_param) OFFSET sqlc.arg(offset_param);

-- name: SearchGlobalNoText :many
SELECT *, COUNT(*) OVER () AS total_count
FROM (
         -- ===============================
         -- 1. Monuments & lieux
         -- ===============================
         SELECT m.id_monument_lieu                           AS id,
                m.titre_monu_lieu                            AS title,
                COALESCE(array_agg(DISTINCT bs.siecle_list) FILTER (WHERE bs.siecle_list IS NOT NULL),
                         '{}')                               AS siecles,
                COALESCE(array_agg(DISTINCT bmn.monu_lieu_nature_type)
                         FILTER (WHERE bmn.monu_lieu_nature_type IS NOT NULL),
                         '{}')                               AS natures,
                COALESCE(
                                jsonb_agg(
                                DISTINCT jsonb_build_object(
                                        'id', tm.id_media,
                                        'titre', tm.titre_media
                                         )
                                         ) FILTER (
                                    WHERE tm.titre_media IS NOT NULL
                                AND tm.titre_media <> ''
                                AND tm.chemin_media IS NOT NULL
                                AND tm.chemin_media <> ''
                                AND jsonb_typeof(tm.chemin_media::jsonb) = 'array'
                                AND EXISTS (SELECT 1
                                            FROM jsonb_array_elements(tm.chemin_media::jsonb) AS elem
                                            WHERE COALESCE(elem ->> 'path', '') <> '')
                                    ),
                                '[]'::jsonb
                )                                            AS medias,
                '{}'::text[]                                 AS professions,
                'monuments_lieux'                            AS source
         FROM t_monuments_lieux m
                  LEFT JOIN loc_communes c ON c.id_commune = m.id_commune
                  LEFT JOIN loc_departements d ON d.id_departement = c.id_departement
                  LEFT JOIN loc_regions r ON r.id_region = d.id_region
                  LEFT JOIN loc_pays p ON p.id_pays = r.id_pays
                  LEFT JOIN cor_siecles_monu_lieu csm ON csm.monument_lieu_id = m.id_monument_lieu
                  LEFT JOIN bib_siecle bs ON bs.id_siecle = csm.siecle_monu_lieu_id
                  LEFT JOIN cor_natures_monu_lieu cnm ON cnm.monument_lieu_id = m.id_monument_lieu
                  LEFT JOIN bib_monu_lieu_natures bmn ON bmn.id_monu_lieu_nature = cnm.monu_lieu_nature_id
                  LEFT JOIN cor_etat_cons_monu_lieu cem ON cem.monument_lieu_id = m.id_monument_lieu
                  LEFT JOIN cor_materiaux_monu_lieu cmm ON cmm.monument_lieu_id = m.id_monument_lieu
                  LEFT JOIN cor_medias_monu_lieu cme ON m.id_monument_lieu = cme.monument_lieu_id
                  LEFT JOIN t_medias tm ON tm.id_media = cme.media_monu_lieu_id
         WHERE sqlc.arg('include_monuments_lieux') = true
           AND ((sqlc.arg('siecles')::int[]) IS NULL OR cardinality(sqlc.arg('siecles')::int[]) = 0 OR
                csm.siecle_monu_lieu_id = ANY (sqlc.arg('siecles')::int[]))
           AND ((sqlc.arg('pays')::int[]) IS NULL OR cardinality(sqlc.arg('pays')::int[]) = 0 OR
                p.id_pays = ANY (sqlc.arg('pays')::int[]))
           AND ((sqlc.arg('region')::int[]) IS NULL OR cardinality(sqlc.arg('region')::int[]) = 0 OR
                r.id_region = ANY (sqlc.arg('region')::int[]))
           AND ((sqlc.arg('departement')::int[]) IS NULL OR cardinality(sqlc.arg('departement')::int[]) = 0 OR
                d.id_departement = ANY (sqlc.arg('departement')::int[]))
           AND ((sqlc.arg('commune')::int[]) IS NULL OR cardinality(sqlc.arg('commune')::int[]) = 0 OR
                c.id_commune = ANY (sqlc.arg('commune')::int[]))
           AND ((sqlc.arg('natures_monu')::int[]) IS NULL OR cardinality(sqlc.arg('natures_monu')::int[]) = 0 OR
                cnm.monu_lieu_nature_id = ANY (sqlc.arg('natures_monu')::int[]))
           AND ((sqlc.arg('etats_monu')::int[]) IS NULL OR cardinality(sqlc.arg('etats_monu')::int[]) = 0 OR
                cnm.monu_lieu_nature_id = ANY (sqlc.arg('etats_monu')::int[]))
           AND ((sqlc.arg('materiaux_monu')::int[]) IS NULL OR cardinality(sqlc.arg('materiaux_monu')::int[]) = 0 OR
                cnm.monu_lieu_nature_id = ANY (sqlc.arg('materiaux_monu')::int[]))
           AND m.publie = true AND m.publication_status = 'PUBLISHED'
         GROUP BY m.id_monument_lieu

         UNION ALL

         -- ===============================
         -- 2. Mobiliers & images
         -- ===============================
         SELECT mob.id_mobilier_image                        AS id,
                mob.titre_mob_img                            AS title,
                COALESCE(array_agg(DISTINCT bs.siecle_list) FILTER (WHERE bs.siecle_list IS NOT NULL),
                         '{}')                               AS siecles,
                COALESCE(array_agg(DISTINCT bmn.nature_type) FILTER (WHERE bmn.nature_type IS NOT NULL),
                         '{}')                               AS natures,
                COALESCE(
                                jsonb_agg(
                                DISTINCT jsonb_build_object(
                                        'id', tm.id_media,
                                        'titre', tm.titre_media
                                         )
                                         ) FILTER (
                                    WHERE tm.titre_media IS NOT NULL
                                AND tm.titre_media <> ''
                                AND tm.chemin_media IS NOT NULL
                                AND tm.chemin_media <> ''
                                AND jsonb_typeof(tm.chemin_media::jsonb) = 'array'
                                AND EXISTS (SELECT 1
                                            FROM jsonb_array_elements(tm.chemin_media::jsonb) AS elem
                                            WHERE COALESCE(elem ->> 'path', '') <> '')
                                    ),
                                '[]'::jsonb
                )                                            AS medias,
                '{}'::text[]                                 AS professions,
                'mobiliers_images'                           AS source
         FROM t_mobiliers_images mob
                  LEFT JOIN loc_communes c ON c.id_commune = mob.id_commune
                  LEFT JOIN loc_departements d ON d.id_departement = c.id_departement
                  LEFT JOIN loc_regions r ON r.id_region = d.id_region
                  LEFT JOIN loc_pays p ON p.id_pays = r.id_pays
                  LEFT JOIN cor_siecles_mob_img csm ON csm.mobilier_image_id = mob.id_mobilier_image
                  LEFT JOIN bib_siecle bs ON bs.id_siecle = csm.siecle_mob_img_id
                  LEFT JOIN cor_natures_mob_img cnm ON cnm.mobilier_image_id = mob.id_mobilier_image
                  LEFT JOIN bib_mob_img_natures bmn ON bmn.id_nature = cnm.nature_id
                  LEFT JOIN cor_etat_cons_mob_img cem ON cem.mobilier_image_id = mob.id_mobilier_image
                  LEFT JOIN cor_materiaux_mob_img cmm ON cmm.mobilier_image_id = mob.id_mobilier_image
                  LEFT JOIN cor_techniques_mob_img ctm ON ctm.mobilier_image_id = mob.id_mobilier_image
                  LEFT JOIN cor_medias_mob_img cme ON mob.id_mobilier_image = cme.mobilier_image_id
                  LEFT JOIN t_medias tm ON tm.id_media = cme.media_mob_img_id
         WHERE sqlc.arg('include_mobiliers_images') = true
           AND ((sqlc.arg('siecles')::int[]) IS NULL OR cardinality(sqlc.arg('siecles')::int[]) = 0 OR
                csm.siecle_mob_img_id = ANY (sqlc.arg('siecles')::int[]))
           AND ((sqlc.arg('pays')::int[]) IS NULL OR cardinality(sqlc.arg('pays')::int[]) = 0 OR
                p.id_pays = ANY (sqlc.arg('pays')::int[]))
           AND ((sqlc.arg('region')::int[]) IS NULL OR cardinality(sqlc.arg('region')::int[]) = 0 OR
                r.id_region = ANY (sqlc.arg('region')::int[]))
           AND ((sqlc.arg('departement')::int[]) IS NULL OR cardinality(sqlc.arg('departement')::int[]) = 0 OR
                d.id_departement = ANY (sqlc.arg('departement')::int[]))
           AND ((sqlc.arg('commune')::int[]) IS NULL OR cardinality(sqlc.arg('commune')::int[]) = 0 OR
                c.id_commune = ANY (sqlc.arg('commune')::int[]))
           AND ((sqlc.arg('natures_mob')::int[]) IS NULL OR cardinality(sqlc.arg('natures_mob')::int[]) = 0 OR
                cnm.nature_id = ANY (sqlc.arg('natures_mob')::int[]))
           AND ((sqlc.arg('etats_mob')::int[]) IS NULL OR cardinality(sqlc.arg('etats_mob')::int[]) = 0 OR
                cem.etat_cons_mob_img_id = ANY (sqlc.arg('etats_mob')::int[]))
           AND ((sqlc.arg('materiaux_mob')::int[]) IS NULL OR cardinality(sqlc.arg('materiaux_mob')::int[]) = 0 OR
                cmm.materiau_mob_img_id = ANY (sqlc.arg('materiaux_mob')::int[]))
           AND ((sqlc.arg('techniques_mob')::int[]) IS NULL OR cardinality(sqlc.arg('techniques_mob')::int[]) = 0 OR
                ctm.technique_id = ANY (sqlc.arg('techniques_mob')::int[]))
           AND mob.publie = true AND mob.publication_status = 'PUBLISHED'
         GROUP BY mob.id_mobilier_image

         UNION ALL

         -- ===============================
         -- 3. Personnes morales
         -- ===============================
         SELECT pm.id_pers_morale                           AS id,
                pm.titre_pers_mo                            AS title,
                COALESCE(array_agg(DISTINCT bs.siecle_list) FILTER (WHERE bs.siecle_list IS NOT NULL),
                         '{}')                              AS siecles,
                COALESCE(array_agg(DISTINCT bpn.pers_mo_nature_type) FILTER (WHERE bpn.pers_mo_nature_type IS NOT NULL),
                         '{}')                              AS natures,
                COALESCE(
                                jsonb_agg(
                                DISTINCT jsonb_build_object(
                                        'id', tm.id_media,
                                        'titre', tm.titre_media
                                         )
                                         ) FILTER (
                                    WHERE tm.titre_media IS NOT NULL
                                AND tm.titre_media <> ''
                                AND tm.chemin_media IS NOT NULL
                                AND tm.chemin_media <> ''
                                AND jsonb_typeof(tm.chemin_media::jsonb) = 'array'
                                AND EXISTS (SELECT 1
                                            FROM jsonb_array_elements(tm.chemin_media::jsonb) AS elem
                                            WHERE COALESCE(elem ->> 'path', '') <> '')
                                    ),
                                '[]'::jsonb
                )                                           AS medias,
                '{}'::text[]                                AS professions,
                'personnes_morales'                         AS source
         FROM t_pers_morales pm
                  LEFT JOIN loc_communes c ON c.id_commune = pm.id_commune
                  LEFT JOIN loc_departements d ON d.id_departement = c.id_departement
                  LEFT JOIN loc_regions r ON r.id_region = d.id_region
                  LEFT JOIN loc_pays p ON p.id_pays = r.id_pays
                  LEFT JOIN cor_siecles_pers_mo csp ON csp.pers_morale_id = pm.id_pers_morale
                  LEFT JOIN bib_siecle bs ON bs.id_siecle = csp.siecle_pers_mo_id
                  LEFT JOIN cor_natures_pers_mo cnp ON cnp.pers_morale_id = pm.id_pers_morale
                  LEFT JOIN bib_pers_mo_natures bpn ON bpn.id_pers_mo_nature = cnp.pers_mo_nature_id
                  LEFT JOIN cor_medias_mob_img cme ON pm.id_pers_morale = cme.mobilier_image_id
                  LEFT JOIN t_medias tm ON tm.id_media = cme.media_mob_img_id
         WHERE sqlc.arg('include_pers_morales') = true
           AND ((sqlc.arg('siecles')::int[]) IS NULL OR cardinality(sqlc.arg('siecles')::int[]) = 0 OR
                csp.siecle_pers_mo_id = ANY (sqlc.arg('siecles')::int[]))
           AND ((sqlc.arg('pays')::int[]) IS NULL OR cardinality(sqlc.arg('pays')::int[]) = 0 OR
                p.id_pays = ANY (sqlc.arg('pays')::int[]))
           AND ((sqlc.arg('region')::int[]) IS NULL OR cardinality(sqlc.arg('region')::int[]) = 0 OR
                r.id_region = ANY (sqlc.arg('region')::int[]))
           AND ((sqlc.arg('departement')::int[]) IS NULL OR cardinality(sqlc.arg('departement')::int[]) = 0 OR
                d.id_departement = ANY (sqlc.arg('departement')::int[]))
           AND ((sqlc.arg('commune')::int[]) IS NULL OR cardinality(sqlc.arg('commune')::int[]) = 0 OR
                c.id_commune = ANY (sqlc.arg('commune')::int[]))
           AND ((sqlc.arg('natures_pers_mo')::int[]) IS NULL OR cardinality(sqlc.arg('natures_pers_mo')::int[]) = 0 OR
                cnp.pers_mo_nature_id = ANY (sqlc.arg('natures_pers_mo')::int[]))
           AND pm.publie = true AND pm.publication_status = 'PUBLISHED'
         GROUP BY pm.id_pers_morale

         UNION ALL

         -- ===============================
         -- 4. Personnes physiques
         -- ===============================
         SELECT pp.id_pers_physique                               AS id,
                pp.prenom_nom_pers_phy                            AS title,
                COALESCE(array_agg(DISTINCT bs.siecle_list) FILTER (WHERE bs.siecle_list IS NOT NULL),
                         '{}')                                    AS siecles,
                '{}'::text[]                                      AS natures,
                COALESCE(
                                jsonb_agg(
                                DISTINCT jsonb_build_object(
                                        'id', tm.id_media,
                                        'titre', tm.titre_media
                                         )
                                         ) FILTER (
                                    WHERE tm.titre_media IS NOT NULL
                                AND tm.titre_media <> ''
                                AND tm.chemin_media IS NOT NULL
                                AND tm.chemin_media <> ''
                                AND jsonb_typeof(tm.chemin_media::jsonb) = 'array'
                                AND EXISTS (SELECT 1
                                            FROM jsonb_array_elements(tm.chemin_media::jsonb) AS elem
                                            WHERE COALESCE(elem ->> 'path', '') <> '')
                                    ),
                                '[]'::jsonb
                )                                                 AS medias,
                COALESCE(array_agg(DISTINCT bpp.profession_type) FILTER (WHERE bpp.profession_type IS NOT NULL),
                         '{}')                                    AS professions,
                'personnes_physiques'                             AS source
         FROM t_pers_physiques pp
                  LEFT JOIN loc_communes c ON c.id_commune = pp.id_commune
                  LEFT JOIN loc_departements d ON d.id_departement = c.id_departement
                  LEFT JOIN loc_regions r ON r.id_region = d.id_region
                  LEFT JOIN loc_pays p ON p.id_pays = r.id_pays
                  LEFT JOIN cor_siecles_pers_phy csp ON csp.pers_physique_id = pp.id_pers_physique
                  LEFT JOIN bib_siecle bs ON bs.id_siecle = csp.siecle_pers_phy_id
                  LEFT JOIN cor_professions_pers_phy cpp ON cpp.pers_physique_id = pp.id_pers_physique
                  LEFT JOIN bib_pers_phy_professions bpp ON bpp.id_profession = cpp.profession_id
                  LEFT JOIN cor_modes_deplacements_pers_phy cmd ON cmd.pers_physique_id = pp.id_pers_physique
                  LEFT JOIN cor_medias_pers_phy cmp ON pp.id_pers_physique = cmp.pers_physique_id
                  LEFT JOIN t_medias tm ON tm.id_media = cmp.media_pers_phy_id
         WHERE sqlc.arg('include_pers_physiques') = true
           AND ((sqlc.arg('siecles')::int[]) IS NULL OR cardinality(sqlc.arg('siecles')::int[]) = 0 OR
                csp.siecle_pers_phy_id = ANY (sqlc.arg('siecles')::int[]))
           AND ((sqlc.arg('pays')::int[]) IS NULL OR cardinality(sqlc.arg('pays')::int[]) = 0 OR
                p.id_pays = ANY (sqlc.arg('pays')::int[]))
           AND ((sqlc.arg('region')::int[]) IS NULL OR cardinality(sqlc.arg('region')::int[]) = 0 OR
                r.id_region = ANY (sqlc.arg('region')::int[]))
           AND ((sqlc.arg('departement')::int[]) IS NULL OR cardinality(sqlc.arg('departement')::int[]) = 0 OR
                d.id_departement = ANY (sqlc.arg('departement')::int[]))
           AND ((sqlc.arg('commune')::int[]) IS NULL OR cardinality(sqlc.arg('commune')::int[]) = 0 OR
                c.id_commune = ANY (sqlc.arg('commune')::int[]))
           AND ((sqlc.arg('professions')::int[]) IS NULL OR cardinality(sqlc.arg('professions')::int[]) = 0 OR
                cpp.profession_id = ANY (sqlc.arg('professions')::int[]))
           AND ((sqlc.arg('modes_deplacements')::int[]) IS NULL OR
                cardinality(sqlc.arg('modes_deplacements')::int[]) = 0 OR
                cmd.mode_deplacement_id = ANY (sqlc.arg('modes_deplacements')::int[]))
           AND pp.publie = true AND pp.publication_status = 'PUBLISHED'
         GROUP BY pp.id_pers_physique) AS results


ORDER BY title
LIMIT sqlc.arg(limit_param) OFFSET sqlc.arg(offset_param);

