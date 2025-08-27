BEGIN;


DROP INDEX IF EXISTS idx_app_users_email;

DROP INDEX IF EXISTS idx_token_password_reset;

DROP TABLE IF EXISTS t_password_resets;

DROP TABLE IF EXISTS t_app_users;

DROP TYPE IF EXISTS user_grade;

DROP TYPE IF EXISTS domaine_expertise;

DROP INDEX IF EXISTS idx_monuments_titre_trgm;
DROP INDEX IF EXISTS idx_mobiliers_titre_trgm;
DROP INDEX IF EXISTS idx_pers_morales_titre_trgm;
DROP INDEX IF EXISTS idx_pers_physiques_nom_trgm;

DELETE FROM t_pers_physiques WHERE publication_status = 'DRAFT';
DELETE FROM t_monuments_lieux WHERE publication_status = 'DRAFT';
DELETE FROM t_mobiliers_images WHERE publication_status = 'DRAFT';
DELETE FROM t_pers_morales WHERE publication_status = 'DRAFT';

DROP INDEX IF EXISTS idx_publication_status_pers_physiques;
DROP INDEX IF EXISTS idx_publication_status_pers_morales;
DROP INDEX IF EXISTS idx_publication_status_mobiliers_images;
DROP INDEX IF EXISTS idx_publication_status_monuments_lieux;

ALTER TABLE t_pers_physiques
    DROP COLUMN IF EXISTS publication_status;
ALTER TABLE t_pers_physiques
    DROP COLUMN IF EXISTS parent_id;

ALTER TABLE t_monuments_lieux
    DROP COLUMN IF EXISTS publication_status;
ALTER TABLE t_monuments_lieux
    DROP COLUMN IF EXISTS parent_id;

ALTER TABLE t_mobiliers_images
    DROP COLUMN IF EXISTS publication_status;
ALTER TABLE t_mobiliers_images
    DROP COLUMN IF EXISTS parent_id;

ALTER TABLE t_pers_morales
    DROP COLUMN IF EXISTS publication_status;
ALTER TABLE t_pers_morales
    DROP COLUMN IF EXISTS parent_id;

DROP TYPE IF EXISTS publication_status;

DROP EXTENSION IF EXISTS pg_trgm;

COMMIT;