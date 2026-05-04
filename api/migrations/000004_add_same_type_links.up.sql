BEGIN;

CREATE TABLE IF NOT EXISTS cor_monu_lieu_monu_lieu (
    monument_lieu_id_1 INTEGER NOT NULL REFERENCES t_monuments_lieux(id_monument_lieu),
    monument_lieu_id_2 INTEGER NOT NULL REFERENCES t_monuments_lieux(id_monument_lieu),
    CHECK (monument_lieu_id_1 <> monument_lieu_id_2)
);

CREATE TABLE IF NOT EXISTS cor_mob_img_mob_img (
    mobilier_image_id_1 INTEGER NOT NULL REFERENCES t_mobiliers_images(id_mobilier_image),
    mobilier_image_id_2 INTEGER NOT NULL REFERENCES t_mobiliers_images(id_mobilier_image),
    CHECK (mobilier_image_id_1 <> mobilier_image_id_2)
);

CREATE TABLE IF NOT EXISTS cor_pers_mo_pers_mo (
    pers_morale_id_1 INTEGER NOT NULL REFERENCES t_pers_morales(id_pers_morale),
    pers_morale_id_2 INTEGER NOT NULL REFERENCES t_pers_morales(id_pers_morale),
    CHECK (pers_morale_id_1 <> pers_morale_id_2)
);

CREATE TABLE IF NOT EXISTS cor_pers_phy_pers_phy (
    pers_physique_id_1 INTEGER NOT NULL REFERENCES t_pers_physiques(id_pers_physique),
    pers_physique_id_2 INTEGER NOT NULL REFERENCES t_pers_physiques(id_pers_physique),
    CHECK (pers_physique_id_1 <> pers_physique_id_2)
);

COMMIT;
