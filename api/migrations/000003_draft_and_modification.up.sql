BEGIN;


ALTER TABLE t_mobiliers_images
    ADD COLUMN user_id TEXT NULL REFERENCES t_app_users(id),
    ADD COLUMN id_departement INTEGER NULL REFERENCES loc_departements(id_departement),
    ADD COLUMN id_region INTEGER NULL REFERENCES loc_regions(id_region);

ALTER TABLE t_monuments_lieux
    ADD COLUMN user_id TEXT NULL REFERENCES t_app_users(id),
    ADD COLUMN id_departement INTEGER NULL REFERENCES loc_departements(id_departement),
    ADD COLUMN id_region INTEGER NULL REFERENCES loc_regions(id_region);

ALTER TABLE t_pers_morales
    ADD COLUMN user_id TEXT NULL REFERENCES t_app_users(id),
    ADD COLUMN id_departement INTEGER NULL REFERENCES loc_departements(id_departement),
    ADD COLUMN id_region INTEGER NULL REFERENCES loc_regions(id_region);

ALTER TABLE t_pers_physiques
    ADD COLUMN user_id TEXT NULL REFERENCES t_app_users(id),
    ADD COLUMN id_departement INTEGER NULL REFERENCES loc_departements(id_departement),
    ADD COLUMN id_region INTEGER NULL REFERENCES loc_regions(id_region);

CREATE TABLE IF NOT EXISTS cor_mob_img_app_users
(
    mobilier_image_id INTEGER NOT NULL,
    user_id           TEXT NOT NULL,
    PRIMARY KEY (mobilier_image_id, user_id),
    FOREIGN KEY (mobilier_image_id) REFERENCES t_mobiliers_images (id_mobilier_image),
    FOREIGN KEY (user_id) REFERENCES t_app_users (id)
);

CREATE TABLE IF NOT EXISTS cor_monu_lieu_app_users
(
    monument_lieu_id INTEGER NOT NULL,
    user_id          TEXT NOT NULL,
    PRIMARY KEY (monument_lieu_id, user_id),
    FOREIGN KEY (monument_lieu_id) REFERENCES t_monuments_lieux (id_monument_lieu),
    FOREIGN KEY (user_id) REFERENCES t_app_users (id)
);

CREATE TABLE IF NOT EXISTS cor_pers_mo_app_users
(
    pers_mo_id INTEGER NOT NULL,
    user_id    TEXT NOT NULL,
    PRIMARY KEY (pers_mo_id, user_id),
    FOREIGN KEY (pers_mo_id) REFERENCES t_pers_morales (id_pers_morale),
    FOREIGN KEY (user_id) REFERENCES t_app_users (id)
);

CREATE TABLE IF NOT EXISTS cor_pers_phy_app_users
(
    pers_phy_id INTEGER NOT NULL,
    user_id     TEXT NOT NULL,
    PRIMARY KEY (pers_phy_id, user_id),
    FOREIGN KEY (pers_phy_id) REFERENCES t_pers_physiques (id_pers_physique),
    FOREIGN KEY (user_id) REFERENCES t_app_users (id)
);

ALTER TABLE bib_auteurs
    ADD COLUMN user_id TEXT NULL REFERENCES t_app_users(id);

ALTER TABLE t_app_events
    ALTER COLUMN user_id DROP NOT NULL;

ALTER TYPE event_type RENAME VALUE 'document_validation' TO 'document_submission_validation';
ALTER TYPE event_type RENAME VALUE 'document_rejection' TO 'document_submission_rejection';
ALTER TYPE event_type ADD VALUE 'document_update_validation';
ALTER TYPE event_type ADD VALUE 'document_update_rejection';

COMMIT;