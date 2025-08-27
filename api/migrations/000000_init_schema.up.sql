--
-- PostgreSQL database dump
--

-- Dumped from database version 14.18 (Homebrew)
-- Dumped by pg_dump version 14.18 (Homebrew)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: xc_au_Lieu d'origine_updated_at(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public."xc_au_Lieu d'origine_updated_at"() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
                          BEGIN
                            NEW."updated_at" = NOW();
                            RETURN NEW;
                          END;
                          $$;


--
-- Name: xc_au_Patmob_updated_at(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public."xc_au_Patmob_updated_at"() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
                          BEGIN
                            NEW."updated_at" = NOW();
                            RETURN NEW;
                          END;
                          $$;


--
-- Name: xc_au_Sheet-1_updated_at(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public."xc_au_Sheet-1_updated_at"() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
                          BEGIN
                            NEW."updated_at" = NOW();
                            RETURN NEW;
                          END;
                          $$;


--
-- Name: xc_au_Table-1_updated_at(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public."xc_au_Table-1_updated_at"() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
                          BEGIN
                            NEW."updated_at" = NOW();
                            RETURN NEW;
                          END;
                          $$;


--
-- Name: xc_au_Table-2_updated_at(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public."xc_au_Table-2_updated_at"() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
                          BEGIN
                            NEW."updated_at" = NOW();
                            RETURN NEW;
                          END;
                          $$;


--
-- Name: xc_au_bib_monu_lieu_natures_csv_updated_at(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.xc_au_bib_monu_lieu_natures_csv_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
                          BEGIN
                            NEW."updated_at" = NOW();
                            RETURN NEW;
                          END;
                          $$;


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: bib_auteurs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.bib_auteurs (
    id_auteur_fiche integer NOT NULL,
    auteur_fiche_nom character varying(255)
);


--
-- Name: bib_etats_conservation; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.bib_etats_conservation (
    id_etat_conservation integer NOT NULL,
    etat_conservation_type character varying(20)
);


--
-- Name: bib_etats_conservation_id_etat_conservation_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.bib_etats_conservation_id_etat_conservation_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: bib_etats_conservation_id_etat_conservation_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.bib_etats_conservation_id_etat_conservation_seq OWNED BY public.bib_etats_conservation.id_etat_conservation;


--
-- Name: bib_materiaux; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.bib_materiaux (
    id_materiau integer NOT NULL,
    materiau_type character varying(255)
);


--
-- Name: bib_materiaux_id_materiau_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.bib_materiaux_id_materiau_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: bib_materiaux_id_materiau_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.bib_materiaux_id_materiau_seq OWNED BY public.bib_materiaux.id_materiau;


--
-- Name: bib_mob_img_natures; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.bib_mob_img_natures (
    id_nature integer NOT NULL,
    nature_type character varying(255)
);


--
-- Name: bib_mob_img_designations_id_designation_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.bib_mob_img_designations_id_designation_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: bib_mob_img_designations_id_designation_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.bib_mob_img_designations_id_designation_seq OWNED BY public.bib_mob_img_natures.id_nature;


--
-- Name: bib_mob_img_techniques; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.bib_mob_img_techniques (
    id_technique integer NOT NULL,
    technique_type character varying(255)
);


--
-- Name: bib_mob_img_techniques_id_technique_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.bib_mob_img_techniques_id_technique_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: bib_mob_img_techniques_id_technique_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.bib_mob_img_techniques_id_technique_seq OWNED BY public.bib_mob_img_techniques.id_technique;


--
-- Name: bib_monu_lieu_natures_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.bib_monu_lieu_natures_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: bib_monu_lieu_natures; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.bib_monu_lieu_natures (
    id_monu_lieu_nature integer DEFAULT nextval('public.bib_monu_lieu_natures_id_seq'::regclass) NOT NULL,
    monu_lieu_nature_type character varying
);


--
-- Name: bib_pers_mo_compostelle_exige; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.bib_pers_mo_compostelle_exige (
    id_compostelle_exige integer NOT NULL,
    compostelle_exige_type character varying(255)
);


--
-- Name: bib_pers_mo_compostelle_exige_id_compostelle_exige_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.bib_pers_mo_compostelle_exige_id_compostelle_exige_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: bib_pers_mo_compostelle_exige_id_compostelle_exige_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.bib_pers_mo_compostelle_exige_id_compostelle_exige_seq OWNED BY public.bib_pers_mo_compostelle_exige.id_compostelle_exige;


--
-- Name: bib_pers_mo_natures; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.bib_pers_mo_natures (
    id_pers_mo_nature integer NOT NULL,
    pers_mo_nature_type character varying(255)
);


--
-- Name: bib_pers_mo_natures_id_pers_mo_nature_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.bib_pers_mo_natures_id_pers_mo_nature_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: bib_pers_mo_natures_id_pers_mo_nature_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.bib_pers_mo_natures_id_pers_mo_nature_seq OWNED BY public.bib_pers_mo_natures.id_pers_mo_nature;


--
-- Name: bib_pers_phy_attestation; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.bib_pers_phy_attestation (
    id_attestation integer NOT NULL,
    attestation_type character varying(255)
);


--
-- Name: bib_pers_phy_attestation_id_attestation_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.bib_pers_phy_attestation_id_attestation_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: bib_pers_phy_attestation_id_attestation_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.bib_pers_phy_attestation_id_attestation_seq OWNED BY public.bib_pers_phy_attestation.id_attestation;


--
-- Name: bib_pers_phy_modes_deplacements; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.bib_pers_phy_modes_deplacements (
    id_mode_deplacement integer NOT NULL,
    mode_deplacement_type character varying(255)
);


--
-- Name: bib_pers_phy_modes_deplacements_id_mode_deplacement_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.bib_pers_phy_modes_deplacements_id_mode_deplacement_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: bib_pers_phy_modes_deplacements_id_mode_deplacement_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.bib_pers_phy_modes_deplacements_id_mode_deplacement_seq OWNED BY public.bib_pers_phy_modes_deplacements.id_mode_deplacement;


--
-- Name: bib_pers_phy_periodes_historiques; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.bib_pers_phy_periodes_historiques (
    id_periode_historique integer NOT NULL,
    periode_historique_type character varying(255)
);


--
-- Name: bib_pers_phy_periodes_historiques_id_periode_historique_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.bib_pers_phy_periodes_historiques_id_periode_historique_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: bib_pers_phy_periodes_historiques_id_periode_historique_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.bib_pers_phy_periodes_historiques_id_periode_historique_seq OWNED BY public.bib_pers_phy_periodes_historiques.id_periode_historique;


--
-- Name: bib_pers_phy_professions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.bib_pers_phy_professions (
    id_profession integer NOT NULL,
    profession_type character varying(255)
);


--
-- Name: bib_pers_phy_professions_id_profession_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.bib_pers_phy_professions_id_profession_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: bib_pers_phy_professions_id_profession_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.bib_pers_phy_professions_id_profession_seq OWNED BY public.bib_pers_phy_professions.id_profession;


--
-- Name: bib_pers_phy_sexes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.bib_pers_phy_sexes (
    id_sexe integer NOT NULL,
    sexe_type character varying(255)
);


--
-- Name: bib_pers_phy_sexes_groupes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.bib_pers_phy_sexes_groupes (
    id_sexe_groupe integer NOT NULL,
    sexe_groupe_type character varying(255)
);


--
-- Name: bib_pers_phy_sexes_groupes_id_sexe_groupe_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.bib_pers_phy_sexes_groupes_id_sexe_groupe_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: bib_pers_phy_sexes_groupes_id_sexe_groupe_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.bib_pers_phy_sexes_groupes_id_sexe_groupe_seq OWNED BY public.bib_pers_phy_sexes_groupes.id_sexe_groupe;


--
-- Name: bib_pers_phy_sexes_id_sexe_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.bib_pers_phy_sexes_id_sexe_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: bib_pers_phy_sexes_id_sexe_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.bib_pers_phy_sexes_id_sexe_seq OWNED BY public.bib_pers_phy_sexes.id_sexe;


--
-- Name: bib_pers_phy_situations_familiales; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.bib_pers_phy_situations_familiales (
    id_situation_familiale integer NOT NULL,
    situation_familiale_type character varying(255)
);


--
-- Name: bib_pers_phy_situations_familiales_id_situation_familiale_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.bib_pers_phy_situations_familiales_id_situation_familiale_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: bib_pers_phy_situations_familiales_id_situation_familiale_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.bib_pers_phy_situations_familiales_id_situation_familiale_seq OWNED BY public.bib_pers_phy_situations_familiales.id_situation_familiale;


--
-- Name: bib_redacteur_id_redacteur_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.bib_redacteur_id_redacteur_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: bib_redacteur_id_redacteur_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.bib_redacteur_id_redacteur_seq OWNED BY public.bib_auteurs.id_auteur_fiche;


--
-- Name: bib_siecle; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.bib_siecle (
    id_siecle integer NOT NULL,
    siecle_list character varying(12)
);


--
-- Name: bib_siecle_id_siecle_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.bib_siecle_id_siecle_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: bib_siecle_id_siecle_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.bib_siecle_id_siecle_seq OWNED BY public.bib_siecle.id_siecle;


--
-- Name: bib_source_auteur; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.bib_source_auteur (
    id_source_auteur integer NOT NULL,
    source_auteur character varying(255)
);


--
-- Name: bib_source_auteur_id_source_auteur_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.bib_source_auteur_id_source_auteur_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: bib_source_auteur_id_source_auteur_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.bib_source_auteur_id_source_auteur_seq OWNED BY public.bib_source_auteur.id_source_auteur;


--
-- Name: bib_source_date; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.bib_source_date (
    id_source_date integer NOT NULL,
    source_date character varying(255)
);


--
-- Name: bib_source_date_id_source_date_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.bib_source_date_id_source_date_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: bib_source_date_id_source_date_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.bib_source_date_id_source_date_seq OWNED BY public.bib_source_date.id_source_date;


--
-- Name: bib_source_type; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.bib_source_type (
    id_source_type integer NOT NULL,
    source_type character varying(255)
);


--
-- Name: bib_source_type_id_source_type_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.bib_source_type_id_source_type_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: bib_source_type_id_source_type_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.bib_source_type_id_source_type_seq OWNED BY public.bib_source_type.id_source_type;


--
-- Name: cor_auteur_fiche_mob_img; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cor_auteur_fiche_mob_img (
    auteur_fiche_mob_img_id integer NOT NULL,
    mobilier_image_id integer NOT NULL
);


--
-- Name: cor_auteur_fiche_monu_lieu; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cor_auteur_fiche_monu_lieu (
    auteur_fiche_monu_lieu_id integer NOT NULL,
    monument_lieu_id integer NOT NULL
);


--
-- Name: cor_auteur_fiche_pers_mo; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cor_auteur_fiche_pers_mo (
    auteur_fiche_pers_mo_id integer NOT NULL,
    pers_morale_id integer NOT NULL
);


--
-- Name: cor_auteur_fiche_pers_phy; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cor_auteur_fiche_pers_phy (
    auteur_fiche_pers_phy_id integer NOT NULL,
    pers_physique_id integer NOT NULL
);


--
-- Name: cor_etat_cons_mob_img; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cor_etat_cons_mob_img (
    etat_cons_mob_img_id integer NOT NULL,
    mobilier_image_id integer NOT NULL
);


--
-- Name: cor_etat_cons_monu_lieu; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cor_etat_cons_monu_lieu (
    etat_cons_monu_lieu_id integer NOT NULL,
    monument_lieu_id integer NOT NULL
);


--
-- Name: cor_materiaux_mob_img; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cor_materiaux_mob_img (
    materiau_mob_img_id integer NOT NULL,
    mobilier_image_id integer NOT NULL
);


--
-- Name: cor_materiaux_monu_lieu; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cor_materiaux_monu_lieu (
    materiau_monu_lieu_id integer NOT NULL,
    monument_lieu_id integer NOT NULL
);


--
-- Name: cor_medias_mob_img; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cor_medias_mob_img (
    media_mob_img_id integer NOT NULL,
    mobilier_image_id integer NOT NULL
);


--
-- Name: cor_medias_monu_lieu; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cor_medias_monu_lieu (
    media_monu_lieu_id integer NOT NULL,
    monument_lieu_id integer NOT NULL
);


--
-- Name: cor_medias_pers_mo; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cor_medias_pers_mo (
    media_pers_mo_id integer NOT NULL,
    pers_morale_id integer NOT NULL
);


--
-- Name: cor_medias_pers_phy; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cor_medias_pers_phy (
    media_pers_phy_id integer NOT NULL,
    pers_physique_id integer NOT NULL
);


--
-- Name: cor_mob_img_pers_mo; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cor_mob_img_pers_mo (
    mobilier_image_id integer NOT NULL,
    pers_morale_id integer NOT NULL
);


--
-- Name: cor_mob_img_pers_phy; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cor_mob_img_pers_phy (
    mobilier_image_id integer NOT NULL,
    pers_physique_id integer NOT NULL
);


--
-- Name: cor_modes_deplacements_pers_phy; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cor_modes_deplacements_pers_phy (
    mode_deplacement_id integer NOT NULL,
    pers_physique_id integer NOT NULL
);


--
-- Name: cor_monu_lieu_mob_img; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cor_monu_lieu_mob_img (
    monument_lieu_id integer NOT NULL,
    mobilier_image_id integer NOT NULL
);


--
-- Name: cor_monu_lieu_pers_mo; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cor_monu_lieu_pers_mo (
    monument_lieu_id integer NOT NULL,
    pers_morale_id integer NOT NULL
);


--
-- Name: cor_monu_lieu_pers_phy; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cor_monu_lieu_pers_phy (
    monu_lieu_id integer NOT NULL,
    pers_phy_id integer NOT NULL
);


--
-- Name: cor_natures_mob_img; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cor_natures_mob_img (
    nature_id integer NOT NULL,
    mobilier_image_id integer NOT NULL
);


--
-- Name: cor_natures_monu_lieu; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cor_natures_monu_lieu (
    monu_lieu_nature_id integer NOT NULL,
    monument_lieu_id integer NOT NULL
);


--
-- Name: cor_natures_pers_mo; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cor_natures_pers_mo (
    pers_mo_nature_id integer NOT NULL,
    pers_morale_id integer NOT NULL
);


--
-- Name: cor_periodes_historiques_pers_phy; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cor_periodes_historiques_pers_phy (
    periode_historique_id integer NOT NULL,
    pers_physique_id integer NOT NULL
);


--
-- Name: cor_pers_phy_pers_mo; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cor_pers_phy_pers_mo (
    pers_physique_id integer NOT NULL,
    pers_morale_id integer NOT NULL
);


--
-- Name: cor_professions_pers_phy; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cor_professions_pers_phy (
    profession_id integer NOT NULL,
    pers_physique_id integer NOT NULL
);


--
-- Name: cor_siecles_mob_img; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cor_siecles_mob_img (
    siecle_mob_img_id integer NOT NULL,
    mobilier_image_id integer NOT NULL
);


--
-- Name: cor_siecles_monu_lieu; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cor_siecles_monu_lieu (
    siecle_monu_lieu_id integer NOT NULL,
    monument_lieu_id integer NOT NULL
);


--
-- Name: cor_siecles_pers_mo; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cor_siecles_pers_mo (
    siecle_pers_mo_id integer NOT NULL,
    pers_morale_id integer NOT NULL
);


--
-- Name: cor_siecles_pers_phy; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cor_siecles_pers_phy (
    siecle_pers_phy_id integer NOT NULL,
    pers_physique_id integer NOT NULL
);


--
-- Name: cor_source_auteur_mob_img; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cor_source_auteur_mob_img (
    source_auteur_mob_img_id integer NOT NULL,
    mobilier_image_id integer NOT NULL
);


--
-- Name: cor_source_auteur_monu_lieu; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cor_source_auteur_monu_lieu (
    source_auteur_monu_lieu_id integer NOT NULL,
    monument_lieu_id integer NOT NULL
);


--
-- Name: cor_source_auteur_pers_mo; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cor_source_auteur_pers_mo (
    source_auteur_pers_mo_id integer NOT NULL,
    pers_morale_id integer NOT NULL
);


--
-- Name: cor_source_auteur_pers_phy; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cor_source_auteur_pers_phy (
    source_auteur_pers_phy_id integer NOT NULL,
    pers_physique_id integer NOT NULL
);


--
-- Name: cor_source_date_mob_img; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cor_source_date_mob_img (
    source_date_mob_img_id integer NOT NULL,
    mobilier_image_id integer NOT NULL
);


--
-- Name: cor_source_date_monu_lieu; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cor_source_date_monu_lieu (
    source_date_monu_lieu_id integer NOT NULL,
    monument_lieu_id integer NOT NULL
);


--
-- Name: cor_source_date_pers_mo; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cor_source_date_pers_mo (
    source_date_pers_mo_id integer NOT NULL,
    pers_morale_id integer NOT NULL
);


--
-- Name: cor_source_date_pers_phy; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cor_source_date_pers_phy (
    source_date_pers_phy_id integer NOT NULL,
    pers_physique_id integer NOT NULL
);


--
-- Name: cor_source_type_mob_img; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cor_source_type_mob_img (
    source_type_mob_img_id integer NOT NULL,
    mobilier_image_id integer NOT NULL
);


--
-- Name: cor_source_type_monu_lieu; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cor_source_type_monu_lieu (
    source_type_monu_lieu_id integer NOT NULL,
    monument_lieu_id integer NOT NULL
);


--
-- Name: cor_source_type_pers_mo; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cor_source_type_pers_mo (
    source_type_pers_mo_id integer NOT NULL,
    pers_morale_id integer NOT NULL
);


--
-- Name: cor_source_type_pers_phy; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cor_source_type_pers_phy (
    source_type_pers_phy_id integer NOT NULL,
    pers_physique_id integer NOT NULL
);


--
-- Name: cor_techniques_mob_img; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cor_techniques_mob_img (
    technique_id integer NOT NULL,
    mobilier_image_id integer NOT NULL
);


--
-- Name: cor_themes_mob_img; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cor_themes_mob_img (
    theme_id integer NOT NULL,
    mob_img_id integer NOT NULL
);


--
-- Name: cor_themes_monu_lieu; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cor_themes_monu_lieu (
    theme_id integer NOT NULL,
    monu_lieu_id integer NOT NULL
);


--
-- Name: cor_themes_pers_mo; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cor_themes_pers_mo (
    theme_id integer NOT NULL,
    pers_mo_id integer NOT NULL
);


--
-- Name: cor_themes_pers_phy; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cor_themes_pers_phy (
    theme_id integer NOT NULL,
    pers_phy_id integer NOT NULL
);


--
-- Name: loc_communes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.loc_communes (
    id_commune integer NOT NULL,
    nom_commune character varying,
    code_postal character varying,
    nom_departement character varying,
    code_commune_insee character varying,
    id_departement integer
);


--
-- Name: loc_commune_transi2_id_commune_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.loc_commune_transi2_id_commune_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: loc_commune_transi2_id_commune_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.loc_commune_transi2_id_commune_seq OWNED BY public.loc_communes.id_commune;


--
-- Name: loc_departements; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.loc_departements (
    id_departement integer NOT NULL,
    nom_departement character varying(255),
    id_region integer
);


--
-- Name: loc_depart_transi_id_departement_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.loc_depart_transi_id_departement_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: loc_depart_transi_id_departement_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.loc_depart_transi_id_departement_seq OWNED BY public.loc_departements.id_departement;


--
-- Name: loc_pays; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.loc_pays (
    id_pays integer NOT NULL,
    nom_pays character varying(255)
);


--
-- Name: loc_pays_id_pays_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.loc_pays_id_pays_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: loc_pays_id_pays_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.loc_pays_id_pays_seq OWNED BY public.loc_pays.id_pays;


--
-- Name: loc_regions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.loc_regions (
    id_region integer NOT NULL,
    nom_region character varying(255),
    id_pays integer
);


--
-- Name: loc_regions_transi_id_region_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.loc_regions_transi_id_region_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: loc_regions_transi_id_region_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.loc_regions_transi_id_region_seq OWNED BY public.loc_regions.id_region;


--
-- Name: nc_evolutions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.nc_evolutions (
    id integer NOT NULL,
    title character varying(255) NOT NULL,
    "titleDown" character varying(255),
    description character varying(255),
    batch integer,
    checksum character varying(255),
    status integer,
    created timestamp with time zone,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


--
-- Name: nc_evolutions_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.nc_evolutions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: nc_evolutions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.nc_evolutions_id_seq OWNED BY public.nc_evolutions.id;


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


--
-- Name: t_medias; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.t_medias (
    id_media integer NOT NULL,
    fiche_associee character varying,
    titre_media character varying,
    "MM11_média_fiche_associée::titre_media" character varying,
    chemin_media text,
    commentaires character varying,
    media_auteur character varying,
    original_connu text,
    original_possede text,
    media_date_de_fabrication character varying,
    date_creation character varying,
    date_maj character varying,
    redacteur character varying,
    contributeur character varying,
    source character varying,
    media_copyright character varying,
    droits_de_diffusion character varying,
    dictionnaire character varying,
    selection character varying,
    adresse_image character varying
);


--
-- Name: t_medias_id_media_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.t_medias_id_media_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: t_medias_id_media_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.t_medias_id_media_seq OWNED BY public.t_medias.id_media;


--
-- Name: t_mobiliers_images_id_mobilier_image_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.t_mobiliers_images_id_mobilier_image_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: t_mobiliers_images; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.t_mobiliers_images (
    id_mobilier_image integer DEFAULT nextval('public.t_mobiliers_images_id_mobilier_image_seq'::regclass) NOT NULL,
    titre_mob_img character varying(255) NOT NULL,
    cote_reference text,
    description text,
    date_fabrication character varying(255),
    auteur_oeuvre text,
    commanditaire text,
    lieu_conservation character varying,
    emplacement text,
    support text,
    proprietaire_actuel text,
    protection_commentaires text,
    dimensions_support text,
    dimensions_image text,
    inscriptions text,
    historique text,
    temoin_commentaires text,
    source text,
    bibliographie text,
    date_maj date,
    schema_descriptif_source character varying(255),
    protection boolean,
    id_filemaker character varying(20),
    id_commune integer,
    publie boolean,
    loc_lieux_dits text,
    id_pays integer,
    tmp_refmedia character varying,
    lieu_origine character varying,
    contributeurs character varying(255),
    "Nature" character varying,
    date_cr_ation date,
    title64 date
);


--
-- Name: t_monuments_lieux_id_monument_lieu_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.t_monuments_lieux_id_monument_lieu_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: t_monuments_lieux; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.t_monuments_lieux (
    id_monument_lieu integer DEFAULT nextval('public.t_monuments_lieux_id_monument_lieu_seq'::regclass) NOT NULL,
    titre_monu_lieu character varying(255) NOT NULL,
    dimensions text,
    altitude character varying(255),
    description text,
    histoire text,
    geolocalisation text,
    emplacement text,
    date_construction character varying(255),
    premiere_mention character varying,
    proprietaire_actuel character varying,
    architecte character varying,
    protection_commentaires text,
    commanditaire character varying,
    source text,
    bibliographie text,
    date_creation date,
    date_maj date,
    protection boolean,
    id_filemaker character varying(20),
    id_commune integer,
    publie boolean,
    loc_lieux_dits character varying,
    id_pays integer,
    tmp_refmedia character varying,
    contributeurs character varying(255),
    date_maj_temp date
);


--
-- Name: t_pers_morales_id_pers_morale_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.t_pers_morales_id_pers_morale_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: t_pers_morales; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.t_pers_morales (
    id_pers_morale integer DEFAULT nextval('public.t_pers_morales_id_pers_morale_seq'::regclass) NOT NULL,
    titre_pers_mo character varying(255),
    simple_mention boolean DEFAULT false,
    objets text,
    creation_date character varying(255),
    date_premiere_mention character varying(255),
    date_derniere_mention character varying(255),
    refondation_date character varying(255),
    date_fin character varying(255),
    texte_statuts boolean,
    origine_sociale_prof character varying(255),
    biens text,
    membres_connus text,
    frequence_reunions text,
    participation_vie_pol text,
    participation_vie_soc text,
    acte_fondation boolean,
    fetes_solennelles text,
    funerailles text,
    autres_fetes text,
    inhumation_costume text,
    fondateurs text,
    statuts text,
    autorisations text,
    historique text,
    sources text,
    bibliographie text,
    commentaires text,
    date_creation date,
    date_maj date,
    id_filemaker character varying(20),
    id_commune integer,
    publie boolean,
    loc_lieux_dits character varying,
    loc_hors_fr character varying,
    id_pays integer,
    tmp_refmedia character varying,
    fonctionnement text,
    contributeurs character varying(255),
    id_compostelle_exige integer
);


--
-- Name: t_pers_physiques_id_pers_physique_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.t_pers_physiques_id_pers_physique_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: t_pers_physiques; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.t_pers_physiques (
    id_pers_physique integer DEFAULT nextval('public.t_pers_physiques_id_pers_physique_seq'::regclass) NOT NULL,
    prenom_nom_pers_phy character varying(255),
    date_pelerinage character varying(255),
    attestation text,
    duree_pelerinage text,
    date_evenement character varying(255),
    nature_evenement text,
    elements_pelerinage text,
    evenements text,
    preparatifs text,
    chemin_suivi text,
    arrivee text,
    retour text,
    non_execution text,
    commutation_voeu text,
    date_naissance character varying(255),
    date_deces character varying(255),
    age character varying(255),
    elements_biographiques text,
    composition_groupe text,
    sources text,
    historiographie text,
    bibliographie text,
    commentaires text,
    date_creation date,
    date_maj date,
    id_compostelle integer,
    id_sexe integer,
    id_situation_familiale integer,
    id_sexe_groupe integer,
    id_filemaker character varying(20),
    id_commune integer,
    publie boolean,
    loc_hors_fr character varying,
    id_pays integer,
    tmp_refmedia character varying,
    contributeurs character varying(255),
    id_attestation integer
);


--
-- Name: t_themes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.t_themes (
    id_theme integer NOT NULL,
    theme_type character varying(255)
);


--
-- Name: t_themes_id_theme_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.t_themes_id_theme_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: t_themes_id_theme_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.t_themes_id_theme_seq OWNED BY public.t_themes.id_theme;


--
-- Name: bib_auteurs id_auteur_fiche; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bib_auteurs ALTER COLUMN id_auteur_fiche SET DEFAULT nextval('public.bib_redacteur_id_redacteur_seq'::regclass);


--
-- Name: bib_etats_conservation id_etat_conservation; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bib_etats_conservation ALTER COLUMN id_etat_conservation SET DEFAULT nextval('public.bib_etats_conservation_id_etat_conservation_seq'::regclass);


--
-- Name: bib_materiaux id_materiau; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bib_materiaux ALTER COLUMN id_materiau SET DEFAULT nextval('public.bib_materiaux_id_materiau_seq'::regclass);


--
-- Name: bib_mob_img_natures id_nature; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bib_mob_img_natures ALTER COLUMN id_nature SET DEFAULT nextval('public.bib_mob_img_designations_id_designation_seq'::regclass);


--
-- Name: bib_mob_img_techniques id_technique; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bib_mob_img_techniques ALTER COLUMN id_technique SET DEFAULT nextval('public.bib_mob_img_techniques_id_technique_seq'::regclass);


--
-- Name: bib_pers_mo_compostelle_exige id_compostelle_exige; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bib_pers_mo_compostelle_exige ALTER COLUMN id_compostelle_exige SET DEFAULT nextval('public.bib_pers_mo_compostelle_exige_id_compostelle_exige_seq'::regclass);


--
-- Name: bib_pers_mo_natures id_pers_mo_nature; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bib_pers_mo_natures ALTER COLUMN id_pers_mo_nature SET DEFAULT nextval('public.bib_pers_mo_natures_id_pers_mo_nature_seq'::regclass);


--
-- Name: bib_pers_phy_attestation id_attestation; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bib_pers_phy_attestation ALTER COLUMN id_attestation SET DEFAULT nextval('public.bib_pers_phy_attestation_id_attestation_seq'::regclass);


--
-- Name: bib_pers_phy_modes_deplacements id_mode_deplacement; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bib_pers_phy_modes_deplacements ALTER COLUMN id_mode_deplacement SET DEFAULT nextval('public.bib_pers_phy_modes_deplacements_id_mode_deplacement_seq'::regclass);


--
-- Name: bib_pers_phy_periodes_historiques id_periode_historique; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bib_pers_phy_periodes_historiques ALTER COLUMN id_periode_historique SET DEFAULT nextval('public.bib_pers_phy_periodes_historiques_id_periode_historique_seq'::regclass);


--
-- Name: bib_pers_phy_professions id_profession; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bib_pers_phy_professions ALTER COLUMN id_profession SET DEFAULT nextval('public.bib_pers_phy_professions_id_profession_seq'::regclass);


--
-- Name: bib_pers_phy_sexes id_sexe; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bib_pers_phy_sexes ALTER COLUMN id_sexe SET DEFAULT nextval('public.bib_pers_phy_sexes_id_sexe_seq'::regclass);


--
-- Name: bib_pers_phy_sexes_groupes id_sexe_groupe; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bib_pers_phy_sexes_groupes ALTER COLUMN id_sexe_groupe SET DEFAULT nextval('public.bib_pers_phy_sexes_groupes_id_sexe_groupe_seq'::regclass);


--
-- Name: bib_pers_phy_situations_familiales id_situation_familiale; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bib_pers_phy_situations_familiales ALTER COLUMN id_situation_familiale SET DEFAULT nextval('public.bib_pers_phy_situations_familiales_id_situation_familiale_seq'::regclass);


--
-- Name: bib_siecle id_siecle; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bib_siecle ALTER COLUMN id_siecle SET DEFAULT nextval('public.bib_siecle_id_siecle_seq'::regclass);


--
-- Name: bib_source_auteur id_source_auteur; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bib_source_auteur ALTER COLUMN id_source_auteur SET DEFAULT nextval('public.bib_source_auteur_id_source_auteur_seq'::regclass);


--
-- Name: bib_source_date id_source_date; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bib_source_date ALTER COLUMN id_source_date SET DEFAULT nextval('public.bib_source_date_id_source_date_seq'::regclass);


--
-- Name: bib_source_type id_source_type; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bib_source_type ALTER COLUMN id_source_type SET DEFAULT nextval('public.bib_source_type_id_source_type_seq'::regclass);


--
-- Name: loc_communes id_commune; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.loc_communes ALTER COLUMN id_commune SET DEFAULT nextval('public.loc_commune_transi2_id_commune_seq'::regclass);


--
-- Name: loc_departements id_departement; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.loc_departements ALTER COLUMN id_departement SET DEFAULT nextval('public.loc_depart_transi_id_departement_seq'::regclass);


--
-- Name: loc_pays id_pays; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.loc_pays ALTER COLUMN id_pays SET DEFAULT nextval('public.loc_pays_id_pays_seq'::regclass);


--
-- Name: loc_regions id_region; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.loc_regions ALTER COLUMN id_region SET DEFAULT nextval('public.loc_regions_transi_id_region_seq'::regclass);


--
-- Name: nc_evolutions id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.nc_evolutions ALTER COLUMN id SET DEFAULT nextval('public.nc_evolutions_id_seq'::regclass);


--
-- Name: t_medias id_media; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.t_medias ALTER COLUMN id_media SET DEFAULT nextval('public.t_medias_id_media_seq'::regclass);


--
-- Name: t_themes id_theme; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.t_themes ALTER COLUMN id_theme SET DEFAULT nextval('public.t_themes_id_theme_seq'::regclass);


--
-- Name: bib_auteurs bib_auteur_fiche_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bib_auteurs
    ADD CONSTRAINT bib_auteur_fiche_pkey PRIMARY KEY (id_auteur_fiche);


--
-- Name: bib_etats_conservation bib_etats_conservation_etat_conservation_type_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bib_etats_conservation
    ADD CONSTRAINT bib_etats_conservation_etat_conservation_type_key UNIQUE (etat_conservation_type);


--
-- Name: bib_etats_conservation bib_etats_conservation_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bib_etats_conservation
    ADD CONSTRAINT bib_etats_conservation_pkey PRIMARY KEY (id_etat_conservation);


--
-- Name: bib_materiaux bib_materiaux_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bib_materiaux
    ADD CONSTRAINT bib_materiaux_pkey PRIMARY KEY (id_materiau);


--
-- Name: bib_mob_img_natures bib_mob_img_natures_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bib_mob_img_natures
    ADD CONSTRAINT bib_mob_img_natures_pkey PRIMARY KEY (id_nature);


--
-- Name: bib_mob_img_techniques bib_mob_img_techniques_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bib_mob_img_techniques
    ADD CONSTRAINT bib_mob_img_techniques_pkey PRIMARY KEY (id_technique);


--
-- Name: bib_monu_lieu_natures bib_monu_lieu_natures_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bib_monu_lieu_natures
    ADD CONSTRAINT bib_monu_lieu_natures_pkey PRIMARY KEY (id_monu_lieu_nature);


--
-- Name: bib_pers_mo_compostelle_exige bib_pers_mo_compostelle_exige_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bib_pers_mo_compostelle_exige
    ADD CONSTRAINT bib_pers_mo_compostelle_exige_pkey PRIMARY KEY (id_compostelle_exige);


--
-- Name: bib_pers_mo_natures bib_pers_mo_natures_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bib_pers_mo_natures
    ADD CONSTRAINT bib_pers_mo_natures_pkey PRIMARY KEY (id_pers_mo_nature);


--
-- Name: bib_pers_phy_attestation bib_pers_phy_attestation_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bib_pers_phy_attestation
    ADD CONSTRAINT bib_pers_phy_attestation_pkey PRIMARY KEY (id_attestation);


--
-- Name: bib_pers_phy_modes_deplacements bib_pers_phy_modes_deplacements_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bib_pers_phy_modes_deplacements
    ADD CONSTRAINT bib_pers_phy_modes_deplacements_pkey PRIMARY KEY (id_mode_deplacement);


--
-- Name: bib_pers_phy_periodes_historiques bib_pers_phy_periodes_historiques_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bib_pers_phy_periodes_historiques
    ADD CONSTRAINT bib_pers_phy_periodes_historiques_pkey PRIMARY KEY (id_periode_historique);


--
-- Name: bib_pers_phy_professions bib_pers_phy_professions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bib_pers_phy_professions
    ADD CONSTRAINT bib_pers_phy_professions_pkey PRIMARY KEY (id_profession);


--
-- Name: bib_pers_phy_sexes_groupes bib_pers_phy_sexes_groupes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bib_pers_phy_sexes_groupes
    ADD CONSTRAINT bib_pers_phy_sexes_groupes_pkey PRIMARY KEY (id_sexe_groupe);


--
-- Name: bib_pers_phy_sexes bib_pers_phy_sexes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bib_pers_phy_sexes
    ADD CONSTRAINT bib_pers_phy_sexes_pkey PRIMARY KEY (id_sexe);


--
-- Name: bib_pers_phy_situations_familiales bib_pers_phy_situations_familiales_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bib_pers_phy_situations_familiales
    ADD CONSTRAINT bib_pers_phy_situations_familiales_pkey PRIMARY KEY (id_situation_familiale);


--
-- Name: bib_siecle bib_siecle_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bib_siecle
    ADD CONSTRAINT bib_siecle_pkey PRIMARY KEY (id_siecle);


--
-- Name: bib_source_auteur bib_source_auteur_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bib_source_auteur
    ADD CONSTRAINT bib_source_auteur_pkey PRIMARY KEY (id_source_auteur);


--
-- Name: bib_source_date bib_source_date_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bib_source_date
    ADD CONSTRAINT bib_source_date_pkey PRIMARY KEY (id_source_date);


--
-- Name: bib_source_type bib_source_type_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bib_source_type
    ADD CONSTRAINT bib_source_type_pkey PRIMARY KEY (id_source_type);


--
-- Name: cor_auteur_fiche_mob_img cor_auteur_fiche_mob_img_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_auteur_fiche_mob_img
    ADD CONSTRAINT cor_auteur_fiche_mob_img_pkey PRIMARY KEY (auteur_fiche_mob_img_id, mobilier_image_id);


--
-- Name: cor_auteur_fiche_monu_lieu cor_auteur_fiche_monu_lieu_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_auteur_fiche_monu_lieu
    ADD CONSTRAINT cor_auteur_fiche_monu_lieu_pkey PRIMARY KEY (auteur_fiche_monu_lieu_id, monument_lieu_id);


--
-- Name: cor_auteur_fiche_pers_mo cor_auteur_fiche_pers_mo_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_auteur_fiche_pers_mo
    ADD CONSTRAINT cor_auteur_fiche_pers_mo_pkey PRIMARY KEY (auteur_fiche_pers_mo_id, pers_morale_id);


--
-- Name: cor_auteur_fiche_pers_phy cor_auteur_fiche_pers_phy_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_auteur_fiche_pers_phy
    ADD CONSTRAINT cor_auteur_fiche_pers_phy_pkey PRIMARY KEY (auteur_fiche_pers_phy_id, pers_physique_id);


--
-- Name: cor_etat_cons_mob_img cor_etat_cons_mob_img_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_etat_cons_mob_img
    ADD CONSTRAINT cor_etat_cons_mob_img_pkey PRIMARY KEY (etat_cons_mob_img_id, mobilier_image_id);


--
-- Name: cor_etat_cons_monu_lieu cor_etat_cons_monu_lieu_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_etat_cons_monu_lieu
    ADD CONSTRAINT cor_etat_cons_monu_lieu_pkey PRIMARY KEY (etat_cons_monu_lieu_id, monument_lieu_id);


--
-- Name: cor_materiaux_mob_img cor_materiaux_mob_img_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_materiaux_mob_img
    ADD CONSTRAINT cor_materiaux_mob_img_pkey PRIMARY KEY (materiau_mob_img_id, mobilier_image_id);


--
-- Name: cor_materiaux_monu_lieu cor_materiaux_monu_lieu_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_materiaux_monu_lieu
    ADD CONSTRAINT cor_materiaux_monu_lieu_pkey PRIMARY KEY (materiau_monu_lieu_id, monument_lieu_id);


--
-- Name: cor_medias_mob_img cor_medias_mob_img_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_medias_mob_img
    ADD CONSTRAINT cor_medias_mob_img_pkey PRIMARY KEY (media_mob_img_id, mobilier_image_id);


--
-- Name: cor_medias_monu_lieu cor_medias_monu_lieu_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_medias_monu_lieu
    ADD CONSTRAINT cor_medias_monu_lieu_pkey PRIMARY KEY (media_monu_lieu_id, monument_lieu_id);


--
-- Name: cor_medias_pers_mo cor_medias_pers_mo_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_medias_pers_mo
    ADD CONSTRAINT cor_medias_pers_mo_pkey PRIMARY KEY (media_pers_mo_id, pers_morale_id);


--
-- Name: cor_medias_pers_phy cor_medias_pers_phy_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_medias_pers_phy
    ADD CONSTRAINT cor_medias_pers_phy_pkey PRIMARY KEY (media_pers_phy_id, pers_physique_id);


--
-- Name: cor_mob_img_pers_mo cor_mob_img_pers_mo_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_mob_img_pers_mo
    ADD CONSTRAINT cor_mob_img_pers_mo_pkey PRIMARY KEY (mobilier_image_id, pers_morale_id);


--
-- Name: cor_mob_img_pers_phy cor_mob_img_pers_phy_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_mob_img_pers_phy
    ADD CONSTRAINT cor_mob_img_pers_phy_pkey PRIMARY KEY (mobilier_image_id, pers_physique_id);


--
-- Name: cor_modes_deplacements_pers_phy cor_modes_deplacements_pers_phy_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_modes_deplacements_pers_phy
    ADD CONSTRAINT cor_modes_deplacements_pers_phy_pkey PRIMARY KEY (mode_deplacement_id, pers_physique_id);


--
-- Name: cor_monu_lieu_mob_img cor_monu_lieu_mob_img_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_monu_lieu_mob_img
    ADD CONSTRAINT cor_monu_lieu_mob_img_pkey PRIMARY KEY (monument_lieu_id, mobilier_image_id);


--
-- Name: cor_monu_lieu_pers_mo cor_monu_lieu_pers_mo_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_monu_lieu_pers_mo
    ADD CONSTRAINT cor_monu_lieu_pers_mo_pkey PRIMARY KEY (monument_lieu_id, pers_morale_id);


--
-- Name: cor_monu_lieu_pers_phy cor_monu_lieu_pers_phy_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_monu_lieu_pers_phy
    ADD CONSTRAINT cor_monu_lieu_pers_phy_pkey PRIMARY KEY (monu_lieu_id, pers_phy_id);


--
-- Name: cor_natures_mob_img cor_natures_mob_img_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_natures_mob_img
    ADD CONSTRAINT cor_natures_mob_img_pkey PRIMARY KEY (nature_id, mobilier_image_id);


--
-- Name: cor_natures_monu_lieu cor_natures_monu_lieu_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_natures_monu_lieu
    ADD CONSTRAINT cor_natures_monu_lieu_pkey PRIMARY KEY (monu_lieu_nature_id, monument_lieu_id);


--
-- Name: cor_natures_pers_mo cor_natures_pers_mo_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_natures_pers_mo
    ADD CONSTRAINT cor_natures_pers_mo_pkey PRIMARY KEY (pers_mo_nature_id, pers_morale_id);


--
-- Name: cor_periodes_historiques_pers_phy cor_periodes_historiques_pers_phy_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_periodes_historiques_pers_phy
    ADD CONSTRAINT cor_periodes_historiques_pers_phy_pkey PRIMARY KEY (periode_historique_id, pers_physique_id);


--
-- Name: cor_pers_phy_pers_mo cor_pers_phy_pers_mo_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_pers_phy_pers_mo
    ADD CONSTRAINT cor_pers_phy_pers_mo_pkey PRIMARY KEY (pers_physique_id, pers_morale_id);


--
-- Name: cor_professions_pers_phy cor_professions_pers_phy_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_professions_pers_phy
    ADD CONSTRAINT cor_professions_pers_phy_pkey PRIMARY KEY (profession_id, pers_physique_id);


--
-- Name: cor_siecles_mob_img cor_siecles_mob_img_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_siecles_mob_img
    ADD CONSTRAINT cor_siecles_mob_img_pkey PRIMARY KEY (siecle_mob_img_id, mobilier_image_id);


--
-- Name: cor_siecles_monu_lieu cor_siecles_monu_lieu_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_siecles_monu_lieu
    ADD CONSTRAINT cor_siecles_monu_lieu_pkey PRIMARY KEY (siecle_monu_lieu_id, monument_lieu_id);


--
-- Name: cor_siecles_pers_mo cor_siecles_pers_mo_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_siecles_pers_mo
    ADD CONSTRAINT cor_siecles_pers_mo_pkey PRIMARY KEY (siecle_pers_mo_id, pers_morale_id);


--
-- Name: cor_siecles_pers_phy cor_siecles_pers_phy_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_siecles_pers_phy
    ADD CONSTRAINT cor_siecles_pers_phy_pkey PRIMARY KEY (siecle_pers_phy_id, pers_physique_id);


--
-- Name: cor_source_auteur_mob_img cor_source_auteur_mob_img_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_source_auteur_mob_img
    ADD CONSTRAINT cor_source_auteur_mob_img_pkey PRIMARY KEY (source_auteur_mob_img_id, mobilier_image_id);


--
-- Name: cor_source_auteur_monu_lieu cor_source_auteur_monu_lieu_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_source_auteur_monu_lieu
    ADD CONSTRAINT cor_source_auteur_monu_lieu_pkey PRIMARY KEY (source_auteur_monu_lieu_id, monument_lieu_id);


--
-- Name: cor_source_auteur_pers_mo cor_source_auteur_pers_mo_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_source_auteur_pers_mo
    ADD CONSTRAINT cor_source_auteur_pers_mo_pkey PRIMARY KEY (source_auteur_pers_mo_id, pers_morale_id);


--
-- Name: cor_source_auteur_pers_phy cor_source_auteur_pers_phy_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_source_auteur_pers_phy
    ADD CONSTRAINT cor_source_auteur_pers_phy_pkey PRIMARY KEY (source_auteur_pers_phy_id, pers_physique_id);


--
-- Name: cor_source_date_mob_img cor_source_date_mob_img_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_source_date_mob_img
    ADD CONSTRAINT cor_source_date_mob_img_pkey PRIMARY KEY (source_date_mob_img_id, mobilier_image_id);


--
-- Name: cor_source_date_monu_lieu cor_source_date_monu_lieu_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_source_date_monu_lieu
    ADD CONSTRAINT cor_source_date_monu_lieu_pkey PRIMARY KEY (source_date_monu_lieu_id, monument_lieu_id);


--
-- Name: cor_source_date_pers_mo cor_source_date_pers_mo_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_source_date_pers_mo
    ADD CONSTRAINT cor_source_date_pers_mo_pkey PRIMARY KEY (source_date_pers_mo_id, pers_morale_id);


--
-- Name: cor_source_date_pers_phy cor_source_date_pers_phy_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_source_date_pers_phy
    ADD CONSTRAINT cor_source_date_pers_phy_pkey PRIMARY KEY (source_date_pers_phy_id, pers_physique_id);


--
-- Name: cor_source_type_mob_img cor_source_type_mob_img_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_source_type_mob_img
    ADD CONSTRAINT cor_source_type_mob_img_pkey PRIMARY KEY (source_type_mob_img_id, mobilier_image_id);


--
-- Name: cor_source_type_monu_lieu cor_source_type_monu_lieu_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_source_type_monu_lieu
    ADD CONSTRAINT cor_source_type_monu_lieu_pkey PRIMARY KEY (source_type_monu_lieu_id, monument_lieu_id);


--
-- Name: cor_source_type_pers_mo cor_source_type_pers_mo_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_source_type_pers_mo
    ADD CONSTRAINT cor_source_type_pers_mo_pkey PRIMARY KEY (source_type_pers_mo_id, pers_morale_id);


--
-- Name: cor_source_type_pers_phy cor_source_type_pers_phy_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_source_type_pers_phy
    ADD CONSTRAINT cor_source_type_pers_phy_pkey PRIMARY KEY (source_type_pers_phy_id, pers_physique_id);


--
-- Name: cor_techniques_mob_img cor_techniques_mob_img_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_techniques_mob_img
    ADD CONSTRAINT cor_techniques_mob_img_pkey PRIMARY KEY (technique_id, mobilier_image_id);


--
-- Name: cor_themes_mob_img cor_themes_mob_img_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_themes_mob_img
    ADD CONSTRAINT cor_themes_mob_img_pkey PRIMARY KEY (theme_id, mob_img_id);


--
-- Name: cor_themes_monu_lieu cor_themes_monu_lieu_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_themes_monu_lieu
    ADD CONSTRAINT cor_themes_monu_lieu_pkey PRIMARY KEY (theme_id, monu_lieu_id);


--
-- Name: cor_themes_pers_mo cor_themes_pers_mo_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_themes_pers_mo
    ADD CONSTRAINT cor_themes_pers_mo_pkey PRIMARY KEY (theme_id, pers_mo_id);


--
-- Name: cor_themes_pers_phy cor_themes_pers_phy_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_themes_pers_phy
    ADD CONSTRAINT cor_themes_pers_phy_pkey PRIMARY KEY (theme_id, pers_phy_id);


--
-- Name: loc_communes loc_commune_transi2_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.loc_communes
    ADD CONSTRAINT loc_commune_transi2_pkey PRIMARY KEY (id_commune);


--
-- Name: loc_departements loc_depart_transi_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.loc_departements
    ADD CONSTRAINT loc_depart_transi_pkey PRIMARY KEY (id_departement);


--
-- Name: loc_pays loc_pays_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.loc_pays
    ADD CONSTRAINT loc_pays_pkey PRIMARY KEY (id_pays);


--
-- Name: loc_regions loc_regions_transi_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.loc_regions
    ADD CONSTRAINT loc_regions_transi_pkey PRIMARY KEY (id_region);


--
-- Name: nc_evolutions nc_evolutions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.nc_evolutions
    ADD CONSTRAINT nc_evolutions_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: t_medias t_medias_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.t_medias
    ADD CONSTRAINT t_medias_pkey PRIMARY KEY (id_media);


--
-- Name: t_mobiliers_images t_mobiliers_images_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.t_mobiliers_images
    ADD CONSTRAINT t_mobiliers_images_pkey PRIMARY KEY (id_mobilier_image);


--
-- Name: t_monuments_lieux t_monuments_lieux_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.t_monuments_lieux
    ADD CONSTRAINT t_monuments_lieux_pkey PRIMARY KEY (id_monument_lieu);


--
-- Name: t_pers_morales t_pers_morales_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.t_pers_morales
    ADD CONSTRAINT t_pers_morales_pkey PRIMARY KEY (id_pers_morale);


--
-- Name: t_pers_physiques t_pers_physiques_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.t_pers_physiques
    ADD CONSTRAINT t_pers_physiques_pkey PRIMARY KEY (id_pers_physique);


--
-- Name: t_themes t_themes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.t_themes
    ADD CONSTRAINT t_themes_pkey PRIMARY KEY (id_theme);


--
-- Name: cor_auteur_fiche_mob_img cor_auteur_fiche_mob_img_auteur_fiche_mob_img_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_auteur_fiche_mob_img
    ADD CONSTRAINT cor_auteur_fiche_mob_img_auteur_fiche_mob_img_id_fkey FOREIGN KEY (auteur_fiche_mob_img_id) REFERENCES public.bib_auteurs(id_auteur_fiche) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_auteur_fiche_mob_img cor_auteur_fiche_mob_img_mobilier_image_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_auteur_fiche_mob_img
    ADD CONSTRAINT cor_auteur_fiche_mob_img_mobilier_image_id_fkey FOREIGN KEY (mobilier_image_id) REFERENCES public.t_mobiliers_images(id_mobilier_image) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_auteur_fiche_monu_lieu cor_auteur_fiche_monu_lieu_auteur_fiche_monu_lieu_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_auteur_fiche_monu_lieu
    ADD CONSTRAINT cor_auteur_fiche_monu_lieu_auteur_fiche_monu_lieu_id_fkey FOREIGN KEY (auteur_fiche_monu_lieu_id) REFERENCES public.bib_auteurs(id_auteur_fiche) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_auteur_fiche_monu_lieu cor_auteur_fiche_monu_lieu_monument_lieu_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_auteur_fiche_monu_lieu
    ADD CONSTRAINT cor_auteur_fiche_monu_lieu_monument_lieu_id_fkey FOREIGN KEY (monument_lieu_id) REFERENCES public.t_monuments_lieux(id_monument_lieu) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_auteur_fiche_pers_mo cor_auteur_fiche_pers_mo_auteur_fiche_pers_mo_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_auteur_fiche_pers_mo
    ADD CONSTRAINT cor_auteur_fiche_pers_mo_auteur_fiche_pers_mo_id_fkey FOREIGN KEY (auteur_fiche_pers_mo_id) REFERENCES public.bib_auteurs(id_auteur_fiche) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_auteur_fiche_pers_mo cor_auteur_fiche_pers_mo_pers_morale_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_auteur_fiche_pers_mo
    ADD CONSTRAINT cor_auteur_fiche_pers_mo_pers_morale_id_fkey FOREIGN KEY (pers_morale_id) REFERENCES public.t_pers_morales(id_pers_morale) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_auteur_fiche_pers_phy cor_auteur_fiche_pers_phy_auteur_fiche_pers_phy_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_auteur_fiche_pers_phy
    ADD CONSTRAINT cor_auteur_fiche_pers_phy_auteur_fiche_pers_phy_id_fkey FOREIGN KEY (auteur_fiche_pers_phy_id) REFERENCES public.bib_auteurs(id_auteur_fiche) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_auteur_fiche_pers_phy cor_auteur_fiche_pers_phy_pers_physique_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_auteur_fiche_pers_phy
    ADD CONSTRAINT cor_auteur_fiche_pers_phy_pers_physique_id_fkey FOREIGN KEY (pers_physique_id) REFERENCES public.t_pers_physiques(id_pers_physique) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_etat_cons_mob_img cor_etat_cons_mob_img_etat_cons_mob_img_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_etat_cons_mob_img
    ADD CONSTRAINT cor_etat_cons_mob_img_etat_cons_mob_img_id_fkey FOREIGN KEY (etat_cons_mob_img_id) REFERENCES public.bib_etats_conservation(id_etat_conservation) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_etat_cons_mob_img cor_etat_cons_mob_img_mobilier_image_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_etat_cons_mob_img
    ADD CONSTRAINT cor_etat_cons_mob_img_mobilier_image_id_fkey FOREIGN KEY (mobilier_image_id) REFERENCES public.t_mobiliers_images(id_mobilier_image) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_etat_cons_monu_lieu cor_etat_cons_monu_lieu_etat_cons_monu_lieu_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_etat_cons_monu_lieu
    ADD CONSTRAINT cor_etat_cons_monu_lieu_etat_cons_monu_lieu_id_fkey FOREIGN KEY (etat_cons_monu_lieu_id) REFERENCES public.bib_etats_conservation(id_etat_conservation) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_etat_cons_monu_lieu cor_etat_cons_monu_lieu_monument_lieu_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_etat_cons_monu_lieu
    ADD CONSTRAINT cor_etat_cons_monu_lieu_monument_lieu_id_fkey FOREIGN KEY (monument_lieu_id) REFERENCES public.t_monuments_lieux(id_monument_lieu) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_materiaux_mob_img cor_materiaux_mob_img_materiau_mob_img_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_materiaux_mob_img
    ADD CONSTRAINT cor_materiaux_mob_img_materiau_mob_img_id_fkey FOREIGN KEY (materiau_mob_img_id) REFERENCES public.bib_materiaux(id_materiau) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_materiaux_mob_img cor_materiaux_mob_img_mobilier_image_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_materiaux_mob_img
    ADD CONSTRAINT cor_materiaux_mob_img_mobilier_image_id_fkey FOREIGN KEY (mobilier_image_id) REFERENCES public.t_mobiliers_images(id_mobilier_image) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_materiaux_monu_lieu cor_materiaux_monu_lieu_materiau_monu_lieu_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_materiaux_monu_lieu
    ADD CONSTRAINT cor_materiaux_monu_lieu_materiau_monu_lieu_id_fkey FOREIGN KEY (materiau_monu_lieu_id) REFERENCES public.bib_materiaux(id_materiau) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_materiaux_monu_lieu cor_materiaux_monu_lieu_monument_lieu_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_materiaux_monu_lieu
    ADD CONSTRAINT cor_materiaux_monu_lieu_monument_lieu_id_fkey FOREIGN KEY (monument_lieu_id) REFERENCES public.t_monuments_lieux(id_monument_lieu) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_medias_mob_img cor_medias_mob_img_media_mob_img_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_medias_mob_img
    ADD CONSTRAINT cor_medias_mob_img_media_mob_img_id_fkey FOREIGN KEY (media_mob_img_id) REFERENCES public.t_medias(id_media) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_medias_mob_img cor_medias_mob_img_mobilier_image_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_medias_mob_img
    ADD CONSTRAINT cor_medias_mob_img_mobilier_image_id_fkey FOREIGN KEY (mobilier_image_id) REFERENCES public.t_mobiliers_images(id_mobilier_image) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_medias_monu_lieu cor_medias_monu_lieu_media_monu_lieu_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_medias_monu_lieu
    ADD CONSTRAINT cor_medias_monu_lieu_media_monu_lieu_id_fkey FOREIGN KEY (media_monu_lieu_id) REFERENCES public.t_medias(id_media) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_medias_monu_lieu cor_medias_monu_lieu_monument_lieu_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_medias_monu_lieu
    ADD CONSTRAINT cor_medias_monu_lieu_monument_lieu_id_fkey FOREIGN KEY (monument_lieu_id) REFERENCES public.t_monuments_lieux(id_monument_lieu) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_medias_pers_mo cor_medias_pers_mo_media_pers_mo_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_medias_pers_mo
    ADD CONSTRAINT cor_medias_pers_mo_media_pers_mo_id_fkey FOREIGN KEY (media_pers_mo_id) REFERENCES public.t_medias(id_media) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_medias_pers_mo cor_medias_pers_mo_pers_morale_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_medias_pers_mo
    ADD CONSTRAINT cor_medias_pers_mo_pers_morale_id_fkey FOREIGN KEY (pers_morale_id) REFERENCES public.t_pers_morales(id_pers_morale) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_medias_pers_phy cor_medias_pers_phy_media_pers_phy_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_medias_pers_phy
    ADD CONSTRAINT cor_medias_pers_phy_media_pers_phy_id_fkey FOREIGN KEY (media_pers_phy_id) REFERENCES public.t_medias(id_media) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_medias_pers_phy cor_medias_pers_phy_pers_physique_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_medias_pers_phy
    ADD CONSTRAINT cor_medias_pers_phy_pers_physique_id_fkey FOREIGN KEY (pers_physique_id) REFERENCES public.t_pers_physiques(id_pers_physique) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_mob_img_pers_mo cor_mob_img_pers_mo_mobilier_image_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_mob_img_pers_mo
    ADD CONSTRAINT cor_mob_img_pers_mo_mobilier_image_id_fkey FOREIGN KEY (mobilier_image_id) REFERENCES public.t_mobiliers_images(id_mobilier_image) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_mob_img_pers_mo cor_mob_img_pers_mo_pers_morale_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_mob_img_pers_mo
    ADD CONSTRAINT cor_mob_img_pers_mo_pers_morale_id_fkey FOREIGN KEY (pers_morale_id) REFERENCES public.t_pers_morales(id_pers_morale) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_mob_img_pers_phy cor_mob_img_pers_phy_mobilier_image_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_mob_img_pers_phy
    ADD CONSTRAINT cor_mob_img_pers_phy_mobilier_image_id_fkey FOREIGN KEY (mobilier_image_id) REFERENCES public.t_mobiliers_images(id_mobilier_image) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_mob_img_pers_phy cor_mob_img_pers_phy_pers_physique_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_mob_img_pers_phy
    ADD CONSTRAINT cor_mob_img_pers_phy_pers_physique_id_fkey FOREIGN KEY (pers_physique_id) REFERENCES public.t_pers_physiques(id_pers_physique) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_modes_deplacements_pers_phy cor_modes_deplacements_pers_phy_mode_deplacement_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_modes_deplacements_pers_phy
    ADD CONSTRAINT cor_modes_deplacements_pers_phy_mode_deplacement_id_fkey FOREIGN KEY (mode_deplacement_id) REFERENCES public.bib_pers_phy_modes_deplacements(id_mode_deplacement) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_modes_deplacements_pers_phy cor_modes_deplacements_pers_phy_pers_physique_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_modes_deplacements_pers_phy
    ADD CONSTRAINT cor_modes_deplacements_pers_phy_pers_physique_id_fkey FOREIGN KEY (pers_physique_id) REFERENCES public.t_pers_physiques(id_pers_physique) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_monu_lieu_mob_img cor_monu_lieu_mob_img_mobilier_image_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_monu_lieu_mob_img
    ADD CONSTRAINT cor_monu_lieu_mob_img_mobilier_image_id_fkey FOREIGN KEY (mobilier_image_id) REFERENCES public.t_mobiliers_images(id_mobilier_image) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_monu_lieu_mob_img cor_monu_lieu_mob_img_monument_lieu_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_monu_lieu_mob_img
    ADD CONSTRAINT cor_monu_lieu_mob_img_monument_lieu_id_fkey FOREIGN KEY (monument_lieu_id) REFERENCES public.t_monuments_lieux(id_monument_lieu) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_monu_lieu_pers_mo cor_monu_lieu_pers_mo_monument_lieu_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_monu_lieu_pers_mo
    ADD CONSTRAINT cor_monu_lieu_pers_mo_monument_lieu_id_fkey FOREIGN KEY (monument_lieu_id) REFERENCES public.t_monuments_lieux(id_monument_lieu) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_monu_lieu_pers_mo cor_monu_lieu_pers_mo_pers_morale_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_monu_lieu_pers_mo
    ADD CONSTRAINT cor_monu_lieu_pers_mo_pers_morale_id_fkey FOREIGN KEY (pers_morale_id) REFERENCES public.t_pers_morales(id_pers_morale) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_monu_lieu_pers_phy cor_monu_lieu_pers_phy_monu_lieu_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_monu_lieu_pers_phy
    ADD CONSTRAINT cor_monu_lieu_pers_phy_monu_lieu_id_fkey FOREIGN KEY (monu_lieu_id) REFERENCES public.t_monuments_lieux(id_monument_lieu) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_monu_lieu_pers_phy cor_monu_lieu_pers_phy_pers_phy_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_monu_lieu_pers_phy
    ADD CONSTRAINT cor_monu_lieu_pers_phy_pers_phy_id_fkey FOREIGN KEY (pers_phy_id) REFERENCES public.t_pers_physiques(id_pers_physique) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_natures_mob_img cor_natures_mob_img_mobilier_image_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_natures_mob_img
    ADD CONSTRAINT cor_natures_mob_img_mobilier_image_id_fkey FOREIGN KEY (mobilier_image_id) REFERENCES public.t_mobiliers_images(id_mobilier_image) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_natures_mob_img cor_natures_mob_img_nature_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_natures_mob_img
    ADD CONSTRAINT cor_natures_mob_img_nature_id_fkey FOREIGN KEY (nature_id) REFERENCES public.bib_mob_img_natures(id_nature) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_natures_monu_lieu cor_natures_monu_lieu_monu_lieu_nature_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_natures_monu_lieu
    ADD CONSTRAINT cor_natures_monu_lieu_monu_lieu_nature_id_fkey FOREIGN KEY (monu_lieu_nature_id) REFERENCES public.bib_monu_lieu_natures(id_monu_lieu_nature) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_natures_monu_lieu cor_natures_monu_lieu_monument_lieu_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_natures_monu_lieu
    ADD CONSTRAINT cor_natures_monu_lieu_monument_lieu_id_fkey FOREIGN KEY (monument_lieu_id) REFERENCES public.t_monuments_lieux(id_monument_lieu) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_natures_pers_mo cor_natures_pers_mo_pers_mo_nature_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_natures_pers_mo
    ADD CONSTRAINT cor_natures_pers_mo_pers_mo_nature_id_fkey FOREIGN KEY (pers_mo_nature_id) REFERENCES public.bib_pers_mo_natures(id_pers_mo_nature) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_natures_pers_mo cor_natures_pers_mo_pers_morale_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_natures_pers_mo
    ADD CONSTRAINT cor_natures_pers_mo_pers_morale_id_fkey FOREIGN KEY (pers_morale_id) REFERENCES public.t_pers_morales(id_pers_morale) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_periodes_historiques_pers_phy cor_periodes_historiques_pers_phy_periode_historique_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_periodes_historiques_pers_phy
    ADD CONSTRAINT cor_periodes_historiques_pers_phy_periode_historique_id_fkey FOREIGN KEY (periode_historique_id) REFERENCES public.bib_pers_phy_periodes_historiques(id_periode_historique) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_periodes_historiques_pers_phy cor_periodes_historiques_pers_phy_pers_physique_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_periodes_historiques_pers_phy
    ADD CONSTRAINT cor_periodes_historiques_pers_phy_pers_physique_id_fkey FOREIGN KEY (pers_physique_id) REFERENCES public.t_pers_physiques(id_pers_physique) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_pers_phy_pers_mo cor_pers_phy_pers_mo_pers_morale_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_pers_phy_pers_mo
    ADD CONSTRAINT cor_pers_phy_pers_mo_pers_morale_id_fkey FOREIGN KEY (pers_morale_id) REFERENCES public.t_pers_morales(id_pers_morale) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_pers_phy_pers_mo cor_pers_phy_pers_mo_pers_physique_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_pers_phy_pers_mo
    ADD CONSTRAINT cor_pers_phy_pers_mo_pers_physique_id_fkey FOREIGN KEY (pers_physique_id) REFERENCES public.t_pers_physiques(id_pers_physique) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_professions_pers_phy cor_professions_pers_phy_pers_physique_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_professions_pers_phy
    ADD CONSTRAINT cor_professions_pers_phy_pers_physique_id_fkey FOREIGN KEY (pers_physique_id) REFERENCES public.t_pers_physiques(id_pers_physique) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_professions_pers_phy cor_professions_pers_phy_profession_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_professions_pers_phy
    ADD CONSTRAINT cor_professions_pers_phy_profession_id_fkey FOREIGN KEY (profession_id) REFERENCES public.bib_pers_phy_professions(id_profession) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_siecles_mob_img cor_siecles_mob_img_mobilier_image_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_siecles_mob_img
    ADD CONSTRAINT cor_siecles_mob_img_mobilier_image_id_fkey FOREIGN KEY (mobilier_image_id) REFERENCES public.t_mobiliers_images(id_mobilier_image) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_siecles_mob_img cor_siecles_mob_img_siecle_mob_img_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_siecles_mob_img
    ADD CONSTRAINT cor_siecles_mob_img_siecle_mob_img_id_fkey FOREIGN KEY (siecle_mob_img_id) REFERENCES public.bib_siecle(id_siecle) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_siecles_monu_lieu cor_siecles_monu_lieu_monument_lieu_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_siecles_monu_lieu
    ADD CONSTRAINT cor_siecles_monu_lieu_monument_lieu_id_fkey FOREIGN KEY (monument_lieu_id) REFERENCES public.t_monuments_lieux(id_monument_lieu) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_siecles_monu_lieu cor_siecles_monu_lieu_siecle_monu_lieu_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_siecles_monu_lieu
    ADD CONSTRAINT cor_siecles_monu_lieu_siecle_monu_lieu_id_fkey FOREIGN KEY (siecle_monu_lieu_id) REFERENCES public.bib_siecle(id_siecle) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_siecles_pers_mo cor_siecles_pers_mo_pers_morale_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_siecles_pers_mo
    ADD CONSTRAINT cor_siecles_pers_mo_pers_morale_id_fkey FOREIGN KEY (pers_morale_id) REFERENCES public.t_pers_morales(id_pers_morale) ON UPDATE CASCADE NOT VALID;


--
-- Name: cor_siecles_pers_mo cor_siecles_pers_mo_siecle_pers_mo_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_siecles_pers_mo
    ADD CONSTRAINT cor_siecles_pers_mo_siecle_pers_mo_id_fkey FOREIGN KEY (siecle_pers_mo_id) REFERENCES public.bib_siecle(id_siecle) ON UPDATE CASCADE NOT VALID;


--
-- Name: cor_siecles_pers_phy cor_siecles_pers_phy_pers_physique_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_siecles_pers_phy
    ADD CONSTRAINT cor_siecles_pers_phy_pers_physique_id_fkey FOREIGN KEY (pers_physique_id) REFERENCES public.t_pers_physiques(id_pers_physique) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_siecles_pers_phy cor_siecles_pers_phy_siecle_pers_phy_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_siecles_pers_phy
    ADD CONSTRAINT cor_siecles_pers_phy_siecle_pers_phy_id_fkey FOREIGN KEY (siecle_pers_phy_id) REFERENCES public.bib_siecle(id_siecle) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_source_auteur_mob_img cor_source_auteur_mob_img_mobilier_image_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_source_auteur_mob_img
    ADD CONSTRAINT cor_source_auteur_mob_img_mobilier_image_id_fkey FOREIGN KEY (mobilier_image_id) REFERENCES public.t_mobiliers_images(id_mobilier_image) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_source_auteur_mob_img cor_source_auteur_mob_img_source_auteur_mob_img_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_source_auteur_mob_img
    ADD CONSTRAINT cor_source_auteur_mob_img_source_auteur_mob_img_id_fkey FOREIGN KEY (source_auteur_mob_img_id) REFERENCES public.bib_source_auteur(id_source_auteur) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_source_auteur_monu_lieu cor_source_auteur_monu_lieu_monument_lieu_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_source_auteur_monu_lieu
    ADD CONSTRAINT cor_source_auteur_monu_lieu_monument_lieu_id_fkey FOREIGN KEY (monument_lieu_id) REFERENCES public.t_monuments_lieux(id_monument_lieu) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_source_auteur_monu_lieu cor_source_auteur_monu_lieu_source_auteur_monu_lieu_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_source_auteur_monu_lieu
    ADD CONSTRAINT cor_source_auteur_monu_lieu_source_auteur_monu_lieu_id_fkey FOREIGN KEY (source_auteur_monu_lieu_id) REFERENCES public.bib_source_auteur(id_source_auteur) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_source_auteur_pers_mo cor_source_auteur_pers_mo_pers_morale_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_source_auteur_pers_mo
    ADD CONSTRAINT cor_source_auteur_pers_mo_pers_morale_id_fkey FOREIGN KEY (pers_morale_id) REFERENCES public.t_pers_morales(id_pers_morale) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_source_auteur_pers_mo cor_source_auteur_pers_mo_source_auteur_pers_mo_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_source_auteur_pers_mo
    ADD CONSTRAINT cor_source_auteur_pers_mo_source_auteur_pers_mo_id_fkey FOREIGN KEY (source_auteur_pers_mo_id) REFERENCES public.bib_source_auteur(id_source_auteur) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_source_auteur_pers_phy cor_source_auteur_pers_phy_pers_physique_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_source_auteur_pers_phy
    ADD CONSTRAINT cor_source_auteur_pers_phy_pers_physique_id_fkey FOREIGN KEY (pers_physique_id) REFERENCES public.t_pers_physiques(id_pers_physique) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_source_auteur_pers_phy cor_source_auteur_pers_phy_source_auteur_pers_phy_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_source_auteur_pers_phy
    ADD CONSTRAINT cor_source_auteur_pers_phy_source_auteur_pers_phy_id_fkey FOREIGN KEY (source_auteur_pers_phy_id) REFERENCES public.bib_source_auteur(id_source_auteur) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_source_date_mob_img cor_source_date_mob_img_mobilier_image_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_source_date_mob_img
    ADD CONSTRAINT cor_source_date_mob_img_mobilier_image_id_fkey FOREIGN KEY (mobilier_image_id) REFERENCES public.t_mobiliers_images(id_mobilier_image) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_source_date_mob_img cor_source_date_mob_img_source_date_mob_img_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_source_date_mob_img
    ADD CONSTRAINT cor_source_date_mob_img_source_date_mob_img_id_fkey FOREIGN KEY (source_date_mob_img_id) REFERENCES public.bib_source_date(id_source_date) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_source_date_monu_lieu cor_source_date_monu_lieu_monument_lieu_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_source_date_monu_lieu
    ADD CONSTRAINT cor_source_date_monu_lieu_monument_lieu_id_fkey FOREIGN KEY (monument_lieu_id) REFERENCES public.t_monuments_lieux(id_monument_lieu) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_source_date_monu_lieu cor_source_date_monu_lieu_source_date_monu_lieu_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_source_date_monu_lieu
    ADD CONSTRAINT cor_source_date_monu_lieu_source_date_monu_lieu_id_fkey FOREIGN KEY (source_date_monu_lieu_id) REFERENCES public.bib_source_date(id_source_date) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_source_date_pers_mo cor_source_date_pers_mo_pers_morale_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_source_date_pers_mo
    ADD CONSTRAINT cor_source_date_pers_mo_pers_morale_id_fkey FOREIGN KEY (pers_morale_id) REFERENCES public.t_pers_morales(id_pers_morale) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_source_date_pers_mo cor_source_date_pers_mo_source_date_pers_mo_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_source_date_pers_mo
    ADD CONSTRAINT cor_source_date_pers_mo_source_date_pers_mo_id_fkey FOREIGN KEY (source_date_pers_mo_id) REFERENCES public.bib_source_date(id_source_date) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_source_date_pers_phy cor_source_date_pers_phy_pers_physique_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_source_date_pers_phy
    ADD CONSTRAINT cor_source_date_pers_phy_pers_physique_id_fkey FOREIGN KEY (pers_physique_id) REFERENCES public.t_pers_physiques(id_pers_physique) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_source_date_pers_phy cor_source_date_pers_phy_source_date_pers_phy_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_source_date_pers_phy
    ADD CONSTRAINT cor_source_date_pers_phy_source_date_pers_phy_id_fkey FOREIGN KEY (source_date_pers_phy_id) REFERENCES public.bib_source_date(id_source_date) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_source_type_mob_img cor_source_type_mob_img_mobilier_image_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_source_type_mob_img
    ADD CONSTRAINT cor_source_type_mob_img_mobilier_image_id_fkey FOREIGN KEY (mobilier_image_id) REFERENCES public.t_mobiliers_images(id_mobilier_image) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_source_type_mob_img cor_source_type_mob_img_source_type_mob_img_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_source_type_mob_img
    ADD CONSTRAINT cor_source_type_mob_img_source_type_mob_img_id_fkey FOREIGN KEY (source_type_mob_img_id) REFERENCES public.bib_source_type(id_source_type) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_source_type_monu_lieu cor_source_type_monu_lieu_monument_lieu_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_source_type_monu_lieu
    ADD CONSTRAINT cor_source_type_monu_lieu_monument_lieu_id_fkey FOREIGN KEY (monument_lieu_id) REFERENCES public.t_monuments_lieux(id_monument_lieu) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_source_type_monu_lieu cor_source_type_monu_lieu_source_type_monu_lieu_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_source_type_monu_lieu
    ADD CONSTRAINT cor_source_type_monu_lieu_source_type_monu_lieu_id_fkey FOREIGN KEY (source_type_monu_lieu_id) REFERENCES public.bib_source_type(id_source_type) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_source_type_pers_mo cor_source_type_pers_mo_pers_morale_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_source_type_pers_mo
    ADD CONSTRAINT cor_source_type_pers_mo_pers_morale_id_fkey FOREIGN KEY (pers_morale_id) REFERENCES public.t_pers_morales(id_pers_morale) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_source_type_pers_mo cor_source_type_pers_mo_source_type_pers_mo_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_source_type_pers_mo
    ADD CONSTRAINT cor_source_type_pers_mo_source_type_pers_mo_id_fkey FOREIGN KEY (source_type_pers_mo_id) REFERENCES public.bib_source_type(id_source_type) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_source_type_pers_phy cor_source_type_pers_phy_pers_physique_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_source_type_pers_phy
    ADD CONSTRAINT cor_source_type_pers_phy_pers_physique_id_fkey FOREIGN KEY (pers_physique_id) REFERENCES public.t_pers_physiques(id_pers_physique) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_source_type_pers_phy cor_source_type_pers_phy_source_type_pers_phy_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_source_type_pers_phy
    ADD CONSTRAINT cor_source_type_pers_phy_source_type_pers_phy_id_fkey FOREIGN KEY (source_type_pers_phy_id) REFERENCES public.bib_source_type(id_source_type) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_techniques_mob_img cor_techniques_mob_img_mobilier_image_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_techniques_mob_img
    ADD CONSTRAINT cor_techniques_mob_img_mobilier_image_id_fkey FOREIGN KEY (mobilier_image_id) REFERENCES public.t_mobiliers_images(id_mobilier_image) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_techniques_mob_img cor_techniques_mob_img_technique_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_techniques_mob_img
    ADD CONSTRAINT cor_techniques_mob_img_technique_id_fkey FOREIGN KEY (technique_id) REFERENCES public.bib_mob_img_techniques(id_technique) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_themes_mob_img cor_themes_mob_img_mob_img_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_themes_mob_img
    ADD CONSTRAINT cor_themes_mob_img_mob_img_id_fkey FOREIGN KEY (mob_img_id) REFERENCES public.t_mobiliers_images(id_mobilier_image) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_themes_mob_img cor_themes_mob_img_theme_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_themes_mob_img
    ADD CONSTRAINT cor_themes_mob_img_theme_id_fkey FOREIGN KEY (theme_id) REFERENCES public.t_themes(id_theme) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_themes_monu_lieu cor_themes_monu_lieu_monu_lieu_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_themes_monu_lieu
    ADD CONSTRAINT cor_themes_monu_lieu_monu_lieu_id_fkey FOREIGN KEY (monu_lieu_id) REFERENCES public.t_monuments_lieux(id_monument_lieu) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_themes_monu_lieu cor_themes_monu_lieu_theme_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_themes_monu_lieu
    ADD CONSTRAINT cor_themes_monu_lieu_theme_id_fkey FOREIGN KEY (theme_id) REFERENCES public.t_themes(id_theme) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_themes_pers_mo cor_themes_pers_mo_pers_mo_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_themes_pers_mo
    ADD CONSTRAINT cor_themes_pers_mo_pers_mo_id_fkey FOREIGN KEY (pers_mo_id) REFERENCES public.t_pers_morales(id_pers_morale) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_themes_pers_mo cor_themes_pers_mo_theme_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_themes_pers_mo
    ADD CONSTRAINT cor_themes_pers_mo_theme_id_fkey FOREIGN KEY (theme_id) REFERENCES public.t_themes(id_theme) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_themes_pers_phy cor_themes_pers_phy_pers_phy_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_themes_pers_phy
    ADD CONSTRAINT cor_themes_pers_phy_pers_phy_id_fkey FOREIGN KEY (pers_phy_id) REFERENCES public.t_pers_physiques(id_pers_physique) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cor_themes_pers_phy cor_themes_pers_phy_theme_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cor_themes_pers_phy
    ADD CONSTRAINT cor_themes_pers_phy_theme_id_fkey FOREIGN KEY (theme_id) REFERENCES public.t_themes(id_theme) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: loc_departements loc_depart_transi_id_region_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.loc_departements
    ADD CONSTRAINT loc_depart_transi_id_region_fkey FOREIGN KEY (id_region) REFERENCES public.loc_regions(id_region) ON UPDATE CASCADE NOT VALID;


--
-- Name: loc_regions loc_regions_transi_id_pays_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.loc_regions
    ADD CONSTRAINT loc_regions_transi_id_pays_fkey FOREIGN KEY (id_pays) REFERENCES public.loc_pays(id_pays) ON UPDATE CASCADE NOT VALID;


--
-- Name: t_mobiliers_images t_mobiliers_images_id_commune_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.t_mobiliers_images
    ADD CONSTRAINT t_mobiliers_images_id_commune_fkey FOREIGN KEY (id_commune) REFERENCES public.loc_communes(id_commune) ON UPDATE CASCADE ON DELETE CASCADE NOT VALID;


--
-- Name: t_mobiliers_images t_mobiliers_images_id_pays_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.t_mobiliers_images
    ADD CONSTRAINT t_mobiliers_images_id_pays_fkey FOREIGN KEY (id_pays) REFERENCES public.loc_pays(id_pays) ON UPDATE CASCADE ON DELETE CASCADE NOT VALID;


--
-- Name: t_monuments_lieux t_monuments_lieux_id_commune_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.t_monuments_lieux
    ADD CONSTRAINT t_monuments_lieux_id_commune_fkey FOREIGN KEY (id_commune) REFERENCES public.loc_communes(id_commune) ON UPDATE CASCADE NOT VALID;


--
-- Name: t_monuments_lieux t_monuments_lieux_id_pays_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.t_monuments_lieux
    ADD CONSTRAINT t_monuments_lieux_id_pays_fkey FOREIGN KEY (id_pays) REFERENCES public.loc_pays(id_pays) ON UPDATE CASCADE NOT VALID;


--
-- Name: t_pers_morales t_pers_morales_id_commune_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.t_pers_morales
    ADD CONSTRAINT t_pers_morales_id_commune_fkey FOREIGN KEY (id_commune) REFERENCES public.loc_communes(id_commune) ON UPDATE CASCADE NOT VALID;


--
-- Name: t_pers_morales t_pers_morales_id_compostelle_exige_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.t_pers_morales
    ADD CONSTRAINT t_pers_morales_id_compostelle_exige_fkey FOREIGN KEY (id_compostelle_exige) REFERENCES public.bib_pers_mo_compostelle_exige(id_compostelle_exige) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: t_pers_morales t_pers_morales_id_pays_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.t_pers_morales
    ADD CONSTRAINT t_pers_morales_id_pays_fkey FOREIGN KEY (id_pays) REFERENCES public.loc_pays(id_pays) ON UPDATE CASCADE NOT VALID;


--
-- Name: t_pers_physiques t_pers_physiques_id_attestation_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.t_pers_physiques
    ADD CONSTRAINT t_pers_physiques_id_attestation_fkey FOREIGN KEY (id_attestation) REFERENCES public.bib_pers_phy_attestation(id_attestation) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: t_pers_physiques t_pers_physiques_id_commune_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.t_pers_physiques
    ADD CONSTRAINT t_pers_physiques_id_commune_fkey FOREIGN KEY (id_commune) REFERENCES public.loc_communes(id_commune) ON UPDATE CASCADE NOT VALID;


--
-- Name: t_pers_physiques t_pers_physiques_id_pays_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.t_pers_physiques
    ADD CONSTRAINT t_pers_physiques_id_pays_fkey FOREIGN KEY (id_pays) REFERENCES public.loc_pays(id_pays) ON UPDATE CASCADE NOT VALID;


--
-- Name: t_pers_physiques t_pers_physiques_id_sexe_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.t_pers_physiques
    ADD CONSTRAINT t_pers_physiques_id_sexe_fkey FOREIGN KEY (id_sexe) REFERENCES public.bib_pers_phy_sexes(id_sexe) ON UPDATE CASCADE NOT VALID;


--
-- Name: t_pers_physiques t_pers_physiques_id_sexe_groupe_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.t_pers_physiques
    ADD CONSTRAINT t_pers_physiques_id_sexe_groupe_fkey FOREIGN KEY (id_sexe_groupe) REFERENCES public.bib_pers_phy_sexes_groupes(id_sexe_groupe) ON UPDATE CASCADE NOT VALID;


--
-- Name: t_pers_physiques t_pers_physiques_id_situation_familiale_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.t_pers_physiques
    ADD CONSTRAINT t_pers_physiques_id_situation_familiale_fkey FOREIGN KEY (id_situation_familiale) REFERENCES public.bib_pers_phy_situations_familiales(id_situation_familiale) ON UPDATE CASCADE NOT VALID;


--
-- PostgreSQL database dump complete
--

