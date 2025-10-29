BEGIN;


DROP TABLE IF EXISTS cor_mob_img_app_users;

DROP TABLE IF EXISTS cor_monu_lieu_app_users;

DROP TABLE IF EXISTS cor_pers_mo_app_users;

DROP TABLE IF EXISTS cor_pers_phy_app_users;

ALTER TABLE bib_auteurs
    DROP COLUMN IF EXISTS user_id;

ALTER TABLE t_mobiliers_images
    DROP COLUMN IF EXISTS id_departement,
    DROP COLUMN IF EXISTS id_region;

ALTER TABLE t_monuments_lieux
    DROP COLUMN IF EXISTS id_departement,
    DROP COLUMN IF EXISTS id_region;

ALTER TABLE t_pers_morales
    DROP COLUMN IF EXISTS id_departement,
    DROP COLUMN IF EXISTS id_region;

ALTER TABLE t_pers_physiques
    DROP COLUMN IF EXISTS id_departement,
    DROP COLUMN IF EXISTS id_region;

COMMIT;