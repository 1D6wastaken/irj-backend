BEGIN;

-- Collation ICU française insensible à la casse ET aux accents (strength = level1).
-- Utilisée dans les ORDER BY sur des colonnes texte pour obtenir un tri "naturel"
-- français (é = e, è = e, ç = c, œ ~ oe, etc.), résolvant le tri Unicode brut qui
-- plaçait "Féricy" après "Frasnes", "Sélestat" après "Suèvres", etc.
CREATE COLLATION IF NOT EXISTS french_ci (
    provider = icu,
    locale = 'fr-u-ks-level1',
    deterministic = false
);

COMMIT;
