BEGIN;

-- Activer l'extension pg_trgm (une seule fois par base)
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Index trigrammes pour accélérer les recherches
CREATE INDEX IF NOT EXISTS idx_monuments_titre_trgm
    ON t_monuments_lieux USING gin (titre_monu_lieu gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_mobiliers_titre_trgm
    ON t_mobiliers_images USING gin (titre_mob_img gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_pers_morales_titre_trgm
    ON t_pers_morales USING gin (titre_pers_mo gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_pers_physiques_nom_trgm
    ON t_pers_physiques USING gin (prenom_nom_pers_phy gin_trgm_ops);

CREATE TYPE domaine_expertise AS ENUM  (
    'ART',
    'ARCHITECTURE',
    'MEDIEVAL',
    'ARCHEOLOGIE',
    'PATRIMOINE',
    'THEOLOGIE',
    'PELERINAGE',
    'AUTRE'
    );

-- Grade utilisateur
CREATE TYPE user_grade AS ENUM (
    'PENDING',
    'ACTIVE',
    'ADMIN'
    );

CREATE TABLE IF NOT EXISTS t_app_users
(
    id            TEXT PRIMARY KEY,
    prenom        TEXT              NOT NULL,
    nom           TEXT              NOT NULL,
    email         TEXT              NOT NULL UNIQUE,
    email_confirm BOOLEAN           NOT NULL DEFAULT FALSE,
    mot_de_passe  TEXT              NOT NULL, -- hashé côté Go
    telephone     VARCHAR(20),
    organisation  TEXT,
    domaine       domaine_expertise NOT NULL,
    motivation    TEXT,
    grade         user_grade        NOT NULL DEFAULT 'PENDING',
    date_creation TIMESTAMPTZ       NOT NULL DEFAULT NOW(),
    date_maj      TIMESTAMPTZ       NOT NULL DEFAULT NOW(),
    last_login    TIMESTAMPTZ
);

CREATE INDEX idx_app_users_email ON t_app_users (email);

CREATE INDEX idx_app_users_grade ON t_app_users (grade);

CREATE TABLE IF NOT EXISTS t_password_resets
(
    user_id    CHAR(36)        NOT NULL REFERENCES t_app_users (id) ON DELETE CASCADE,
    token      CHAR(36) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ     NOT NULL DEFAULT NOW(),

    PRIMARY KEY (user_id)
);

CREATE INDEX idx_token_password_reset ON t_password_resets (token);

CREATE TYPE publication_status AS ENUM (
    'DRAFT',
    'PENDING',
    'PUBLISHED'
);

ALTER TABLE t_monuments_lieux
ADD COLUMN IF NOT EXISTS publication_status publication_status NOT NULL DEFAULT 'DRAFT';

CREATE INDEX idx_publication_status_monuments_lieux ON t_monuments_lieux (publication_status);

UPDATE t_monuments_lieux
SET publication_status = 'PUBLISHED';

ALTER TABLE t_monuments_lieux
ADD COLUMN IF NOT EXISTS parent_id INTEGER NULL REFERENCES t_monuments_lieux(id_monument_lieu);

ALTER TABLE t_mobiliers_images
ADD COLUMN IF NOT EXISTS publication_status publication_status NOT NULL DEFAULT 'DRAFT';

UPDATE t_mobiliers_images
SET publication_status = 'PUBLISHED';

CREATE INDEX idx_publication_status_mobiliers_images ON t_mobiliers_images (publication_status);

ALTER TABLE t_mobiliers_images
ADD COLUMN IF NOT EXISTS parent_id INTEGER NULL REFERENCES t_mobiliers_images(id_mobilier_image);

ALTER TABLE t_pers_morales
ADD COLUMN IF NOT EXISTS publication_status publication_status NOT NULL DEFAULT 'DRAFT';

UPDATE t_pers_morales
SET publication_status = 'PUBLISHED';

CREATE INDEX idx_publication_status_pers_morales ON t_pers_morales (publication_status);

ALTER TABLE t_pers_morales
ADD COLUMN IF NOT EXISTS parent_id INTEGER NULL REFERENCES t_pers_morales(id_pers_morale);

ALTER TABLE t_pers_physiques
ADD COLUMN IF NOT EXISTS publication_status publication_status NOT NULL DEFAULT 'DRAFT';

UPDATE t_pers_physiques
SET publication_status = 'PUBLISHED';

CREATE INDEX idx_publication_status_pers_physiques ON t_pers_physiques (publication_status);

ALTER TABLE t_pers_physiques
ADD COLUMN IF NOT EXISTS parent_id INTEGER NULL REFERENCES t_pers_physiques(id_pers_physique);

COMMIT;