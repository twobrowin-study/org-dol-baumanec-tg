--
-- PostgreSQL database dump
--

-- Dumped from database version 10.20 (Debian 10.20-1.pgdg90+1)
-- Dumped by pg_dump version 10.20 (Debian 10.20-1.pgdg90+1)

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
-- Name: postgres; Type: DATABASE; Schema: -; Owner: postgres
--

CREATE DATABASE postgres WITH TEMPLATE = template0 ENCODING = 'UTF8' LC_COLLATE = 'en_US.utf8' LC_CTYPE = 'en_US.utf8';


ALTER DATABASE postgres OWNER TO postgres;

\connect postgres

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
-- Name: DATABASE postgres; Type: COMMENT; Schema: -; Owner: postgres
--

COMMENT ON DATABASE postgres IS 'default administrative connection database';


--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


--
-- Name: doctor_type; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.doctor_type AS ENUM (
    'analysis',
    'doctor',
    'therapist'
);


ALTER TYPE public.doctor_type OWNER TO postgres;

--
-- Name: sex; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.sex AS ENUM (
    'male',
    'female'
);


ALTER TYPE public.sex OWNER TO postgres;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: doctor; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.doctor (
    id integer NOT NULL,
    title character varying NOT NULL,
    description character varying NOT NULL,
    doctor_type public.doctor_type NOT NULL,
    sex public.sex,
    is_active boolean DEFAULT true NOT NULL,
    is_fun_male boolean DEFAULT false NOT NULL,
    is_fun_female boolean DEFAULT false
);


ALTER TABLE public.doctor OWNER TO postgres;

--
-- Name: user; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."user" (
    id integer NOT NULL,
    chat_id character varying NOT NULL,
    first_name character varying,
    last_name character varying,
    sex public.sex,
    is_active boolean DEFAULT true NOT NULL
);


ALTER TABLE public."user" OWNER TO postgres;

--
-- Name: user_active; Type: VIEW; Schema: public; Owner: postgres
--

CREATE VIEW public.user_active AS
 WITH RECURSIVE user_active(id, chat_id, first_name, last_name, sex) AS (
         SELECT "user".id,
            "user".chat_id,
            "user".first_name,
            "user".last_name,
            "user".sex
           FROM public."user"
          WHERE (("user".first_name IS NOT NULL) AND ("user".last_name IS NOT NULL) AND ("user".sex IS NOT NULL) AND ("user".is_active = true))
        )
 SELECT user_active.id,
    user_active.chat_id,
    user_active.first_name,
    user_active.last_name,
    user_active.sex
   FROM user_active;


ALTER TABLE public.user_active OWNER TO postgres;

--
-- Name: doctor_active_fun; Type: VIEW; Schema: public; Owner: postgres
--

CREATE VIEW public.doctor_active_fun AS
 WITH RECURSIVE doctor_active_fun(id, user_id, title, doctor_type, is_fun) AS (
         SELECT d.id,
            ua.id AS user_id,
            d.title,
            d.doctor_type,
            false AS is_fun
           FROM public.doctor d,
            public.user_active ua
          WHERE ((d.is_active = true) AND (d.is_fun_male = false) AND (d.is_fun_female = false))
        UNION
         SELECT d.id,
            ua.id AS user_id,
            d.title,
            d.doctor_type,
            true AS is_fun
           FROM public.doctor d,
            public.user_active ua
          WHERE ((d.is_active = true) AND (d.is_fun_male = true) AND (ua.sex = 'male'::public.sex))
        UNION
         SELECT d.id,
            ua.id AS user_id,
            d.title,
            d.doctor_type,
            true AS is_fun
           FROM public.doctor d,
            public.user_active ua
          WHERE ((d.is_active = true) AND (d.is_fun_female = true) AND (ua.sex = 'female'::public.sex))
        )
 SELECT doctor_active_fun.id,
    doctor_active_fun.user_id,
    doctor_active_fun.title,
    doctor_active_fun.doctor_type,
    doctor_active_fun.is_fun
   FROM doctor_active_fun;


ALTER TABLE public.doctor_active_fun OWNER TO postgres;

--
-- Name: newtable_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.newtable_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.newtable_id_seq OWNER TO postgres;

--
-- Name: newtable_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.newtable_id_seq OWNED BY public.doctor.id;


--
-- Name: sticker; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.sticker (
    id integer NOT NULL,
    name character varying NOT NULL,
    file character varying NOT NULL,
    is_active boolean DEFAULT true NOT NULL
);


ALTER TABLE public.sticker OWNER TO postgres;

--
-- Name: sticker_active; Type: VIEW; Schema: public; Owner: postgres
--

CREATE VIEW public.sticker_active AS
 WITH RECURSIVE active_sticker(name, file) AS (
         SELECT sticker.name,
            sticker.file
           FROM public.sticker
          WHERE (sticker.is_active = true)
        )
 SELECT active_sticker.name,
    active_sticker.file
   FROM active_sticker;


ALTER TABLE public.sticker_active OWNER TO postgres;

--
-- Name: stiker_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.stiker_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.stiker_id_seq OWNER TO postgres;

--
-- Name: stiker_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.stiker_id_seq OWNED BY public.sticker.id;


--
-- Name: telegram; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.telegram (
    id integer NOT NULL,
    token character varying NOT NULL,
    admin character varying NOT NULL,
    is_active boolean DEFAULT true NOT NULL
);


ALTER TABLE public.telegram OWNER TO postgres;

--
-- Name: telegram_active_single; Type: VIEW; Schema: public; Owner: postgres
--

CREATE VIEW public.telegram_active_single AS
 WITH RECURSIVE telegram_active_single(token, admin) AS (
         SELECT telegram.token,
            telegram.admin
           FROM public.telegram
          WHERE (telegram.is_active = true)
         LIMIT 1
        )
 SELECT telegram_active_single.token,
    telegram_active_single.admin
   FROM telegram_active_single;


ALTER TABLE public.telegram_active_single OWNER TO postgres;

--
-- Name: telegram_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.telegram_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.telegram_id_seq OWNER TO postgres;

--
-- Name: telegram_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.telegram_id_seq OWNED BY public.telegram.id;


--
-- Name: user_doctor_done; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_doctor_done (
    id integer NOT NULL,
    user_id bigint NOT NULL,
    doctor_id bigint NOT NULL,
    "timestamp" timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    is_active boolean DEFAULT true NOT NULL
);


ALTER TABLE public.user_doctor_done OWNER TO postgres;

--
-- Name: user_analysis_to_go; Type: VIEW; Schema: public; Owner: postgres
--

CREATE VIEW public.user_analysis_to_go AS
 WITH RECURSIVE user_analysis_to_go(id, chat_id, title, ready) AS (
         SELECT d.id,
            u.chat_id,
            d.title,
            true AS ready
           FROM public."user" u,
            public.doctor d,
            public.user_doctor_done udd
          WHERE (udd.is_active AND u.is_active AND d.is_active AND (u.id = udd.user_id) AND (d.id = udd.doctor_id) AND ((d.sex IS NULL) OR (d.sex = u.sex)) AND (d.doctor_type = 'analysis'::public.doctor_type))
        UNION
         SELECT d.id,
            u.chat_id,
            d.title,
            false AS ready
           FROM public."user" u,
            public.doctor d
          WHERE (u.is_active AND d.is_active AND ((d.sex IS NULL) OR (d.sex = u.sex)) AND (NOT (d.id IN ( SELECT udd.doctor_id
                   FROM public.user_doctor_done udd
                  WHERE ((udd.user_id = u.id) AND udd.is_active)))) AND (d.doctor_type = 'analysis'::public.doctor_type))
        )
 SELECT user_analysis_to_go.id,
    user_analysis_to_go.chat_id,
    user_analysis_to_go.title,
    user_analysis_to_go.ready
   FROM user_analysis_to_go;


ALTER TABLE public.user_analysis_to_go OWNER TO postgres;

--
-- Name: user_doctor_done_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.user_doctor_done_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.user_doctor_done_id_seq OWNER TO postgres;

--
-- Name: user_doctor_done_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.user_doctor_done_id_seq OWNED BY public.user_doctor_done.id;


--
-- Name: user_doctor_to_go; Type: VIEW; Schema: public; Owner: postgres
--

CREATE VIEW public.user_doctor_to_go AS
 WITH RECURSIVE user_doctor_to_go(id, chat_id, title, ready) AS (
         SELECT d.id,
            u.chat_id,
            d.title,
            true AS ready
           FROM public."user" u,
            public.doctor d,
            public.user_doctor_done udd
          WHERE (udd.is_active AND u.is_active AND d.is_active AND (u.id = udd.user_id) AND (d.id = udd.doctor_id) AND ((d.sex IS NULL) OR (d.sex = u.sex)) AND (d.doctor_type = 'doctor'::public.doctor_type))
        UNION
         SELECT d.id,
            u.chat_id,
            d.title,
            false AS ready
           FROM public."user" u,
            public.doctor d
          WHERE (u.is_active AND d.is_active AND ((d.sex IS NULL) OR (d.sex = u.sex)) AND (NOT (d.id IN ( SELECT udd.doctor_id
                   FROM public.user_doctor_done udd
                  WHERE ((udd.user_id = u.id) AND udd.is_active)))) AND (d.doctor_type = 'doctor'::public.doctor_type))
        )
 SELECT user_doctor_to_go.id,
    user_doctor_to_go.chat_id,
    user_doctor_to_go.title,
    user_doctor_to_go.ready
   FROM user_doctor_to_go;


ALTER TABLE public.user_doctor_to_go OWNER TO postgres;

--
-- Name: user_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.user_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.user_id_seq OWNER TO postgres;

--
-- Name: user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.user_id_seq OWNED BY public."user".id;


--
-- Name: user_therapist_show; Type: VIEW; Schema: public; Owner: postgres
--

CREATE VIEW public.user_therapist_show AS
SELECT
    NULL::integer AS id,
    NULL::character varying AS chat_id,
    NULL::character varying AS title,
    NULL::boolean AS show,
    NULL::boolean AS is_active;


ALTER TABLE public.user_therapist_show OWNER TO postgres;

--
-- Name: user_therapist_to_go; Type: VIEW; Schema: public; Owner: postgres
--

CREATE VIEW public.user_therapist_to_go AS
 WITH RECURSIVE user_therapist_to_go(id, chat_id, title, show, ready) AS (
         SELECT uts.id,
            u.chat_id,
            uts.title,
            uts.show,
            true AS ready
           FROM public."user" u,
            public.user_therapist_show uts,
            public.user_doctor_done udd
          WHERE (udd.is_active AND u.is_active AND uts.is_active AND (u.id = udd.user_id) AND (uts.id = udd.doctor_id) AND ((uts.chat_id)::text = (u.chat_id)::text))
        UNION
         SELECT uts.id,
            u.chat_id,
            uts.title,
            uts.show,
            false AS ready
           FROM public."user" u,
            public.user_therapist_show uts
          WHERE (u.is_active AND uts.is_active AND (NOT (uts.id IN ( SELECT udd.doctor_id
                   FROM public.user_doctor_done udd
                  WHERE (udd.is_active AND (udd.user_id = u.id))))) AND ((uts.chat_id)::text = (u.chat_id)::text))
        )
 SELECT user_therapist_to_go.id,
    user_therapist_to_go.chat_id,
    user_therapist_to_go.title,
    user_therapist_to_go.show,
    user_therapist_to_go.ready
   FROM user_therapist_to_go;


ALTER TABLE public.user_therapist_to_go OWNER TO postgres;

--
-- Name: user_outdata; Type: VIEW; Schema: public; Owner: postgres
--

CREATE VIEW public.user_outdata AS
 WITH RECURSIVE user_outdata(first_name, last_name, sex, title, ready) AS (
         SELECT u.first_name,
            u.last_name,
            u.sex,
            utg.title,
            utg.ready
           FROM public.user_active u,
            ( SELECT user_doctor_to_go.id,
                    user_doctor_to_go.chat_id,
                    user_doctor_to_go.title,
                    user_doctor_to_go.ready
                   FROM public.user_doctor_to_go
                UNION
                 SELECT user_analysis_to_go.id,
                    user_analysis_to_go.chat_id,
                    user_analysis_to_go.title,
                    user_analysis_to_go.ready
                   FROM public.user_analysis_to_go
                UNION
                 SELECT user_therapist_to_go.id,
                    user_therapist_to_go.chat_id,
                    user_therapist_to_go.title,
                    user_therapist_to_go.ready
                   FROM public.user_therapist_to_go) utg
          WHERE ((utg.chat_id)::text = (u.chat_id)::text)
          ORDER BY u.id, utg.id
        )
 SELECT user_outdata.first_name,
    user_outdata.last_name,
    user_outdata.sex,
    user_outdata.title,
    user_outdata.ready
   FROM user_outdata;


ALTER TABLE public.user_outdata OWNER TO postgres;

--
-- Name: user_waiting_anything; Type: VIEW; Schema: public; Owner: postgres
--

CREATE VIEW public.user_waiting_anything AS
 WITH RECURSIVE user_waiting_anything(id, chat_id) AS (
         SELECT "user".id,
            "user".chat_id
           FROM public."user"
          WHERE ((("user".first_name IS NOT NULL) AND ("user".last_name IS NULL)) OR (("user".first_name IS NULL) AND ("user".last_name IS NOT NULL)) OR ("user".is_active = false))
        )
 SELECT user_waiting_anything.id,
    user_waiting_anything.chat_id
   FROM user_waiting_anything;


ALTER TABLE public.user_waiting_anything OWNER TO postgres;

--
-- Name: user_waiting_for_name_and_sex; Type: VIEW; Schema: public; Owner: postgres
--

CREATE VIEW public.user_waiting_for_name_and_sex AS
 WITH RECURSIVE user_waiting_for_name_and_sex(id, chat_id) AS (
         SELECT "user".id,
            "user".chat_id
           FROM public."user"
          WHERE (("user".first_name IS NULL) AND ("user".last_name IS NULL) AND ("user".sex IS NULL))
        )
 SELECT user_waiting_for_name_and_sex.id,
    user_waiting_for_name_and_sex.chat_id
   FROM user_waiting_for_name_and_sex;


ALTER TABLE public.user_waiting_for_name_and_sex OWNER TO postgres;

--
-- Name: user_waiting_for_sex_but_not_name; Type: VIEW; Schema: public; Owner: postgres
--

CREATE VIEW public.user_waiting_for_sex_but_not_name AS
 WITH RECURSIVE user_waiting_for_sex_but_not_name(id, chat_id) AS (
         SELECT "user".id,
            "user".chat_id
           FROM public."user"
          WHERE (("user".first_name IS NOT NULL) AND ("user".last_name IS NOT NULL) AND ("user".sex IS NULL))
        )
 SELECT user_waiting_for_sex_but_not_name.id,
    user_waiting_for_sex_but_not_name.chat_id
   FROM user_waiting_for_sex_but_not_name;


ALTER TABLE public.user_waiting_for_sex_but_not_name OWNER TO postgres;

--
-- Name: user_with_errors; Type: VIEW; Schema: public; Owner: postgres
--

CREATE VIEW public.user_with_errors AS
 WITH RECURSIVE user_waiting_anything(id, chat_id) AS (
         SELECT "user".id,
            "user".chat_id
           FROM public."user"
          WHERE ((("user".first_name IS NULL) AND ("user".last_name IS NOT NULL)) OR (("user".first_name IS NOT NULL) AND ("user".last_name IS NULL)) OR ("user".is_active = false))
        )
 SELECT user_waiting_anything.id,
    user_waiting_anything.chat_id
   FROM user_waiting_anything;


ALTER TABLE public.user_with_errors OWNER TO postgres;

--
-- Name: doctor id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.doctor ALTER COLUMN id SET DEFAULT nextval('public.newtable_id_seq'::regclass);


--
-- Name: sticker id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sticker ALTER COLUMN id SET DEFAULT nextval('public.stiker_id_seq'::regclass);


--
-- Name: telegram id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.telegram ALTER COLUMN id SET DEFAULT nextval('public.telegram_id_seq'::regclass);


--
-- Name: user id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."user" ALTER COLUMN id SET DEFAULT nextval('public.user_id_seq'::regclass);


--
-- Name: user_doctor_done id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_doctor_done ALTER COLUMN id SET DEFAULT nextval('public.user_doctor_done_id_seq'::regclass);


--
-- Name: 1161361; Type: BLOB; Schema: -; Owner: postgres
--

SELECT pg_catalog.lo_create('1161361');


ALTER LARGE OBJECT 1161361 OWNER TO postgres;

--
-- Name: 1429018; Type: BLOB; Schema: -; Owner: postgres
--

SELECT pg_catalog.lo_create('1429018');


ALTER LARGE OBJECT 1429018 OWNER TO postgres;

--
-- Name: 2805874; Type: BLOB; Schema: -; Owner: postgres
--

SELECT pg_catalog.lo_create('2805874');


ALTER LARGE OBJECT 2805874 OWNER TO postgres;

--
-- Name: 4537866; Type: BLOB; Schema: -; Owner: postgres
--

SELECT pg_catalog.lo_create('4537866');


ALTER LARGE OBJECT 4537866 OWNER TO postgres;

--
-- Name: 5786037; Type: BLOB; Schema: -; Owner: postgres
--

SELECT pg_catalog.lo_create('5786037');


ALTER LARGE OBJECT 5786037 OWNER TO postgres;

--
-- Name: 6413713; Type: BLOB; Schema: -; Owner: postgres
--

SELECT pg_catalog.lo_create('6413713');


ALTER LARGE OBJECT 6413713 OWNER TO postgres;

--
-- Data for Name: doctor; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.doctor (id, title, description, doctor_type, sex, is_active, is_fun_male, is_fun_female) VALUES (1, 'Отоларинголог каб. 4', 'Врач - Отоларинголог', 'doctor', NULL, true, false, false);
INSERT INTO public.doctor (id, title, description, doctor_type, sex, is_active, is_fun_male, is_fun_female) VALUES (3, 'Дерматовенеролог каб. 3', 'Врач - Дерматовенеролог', 'doctor', NULL, true, false, false);
INSERT INTO public.doctor (id, title, description, doctor_type, sex, is_active, is_fun_male, is_fun_female) VALUES (6, 'Невролог каб. 204 эт. 2', 'Врач - Невролог', 'doctor', NULL, true, false, false);
INSERT INTO public.doctor (id, title, description, doctor_type, sex, is_active, is_fun_male, is_fun_female) VALUES (7, 'Акушер-гинеколог каб. 201 эт. 2', 'Врач - Акушер-гинеколог', 'doctor', 'female', true, false, false);
INSERT INTO public.doctor (id, title, description, doctor_type, sex, is_active, is_fun_male, is_fun_female) VALUES (8, 'Стоматолог эт. 3', 'Врач - Стоматолог', 'doctor', NULL, true, false, false);
INSERT INTO public.doctor (id, title, description, doctor_type, sex, is_active, is_fun_male, is_fun_female) VALUES (24, 'Анализ крови и мазки зева, носа, гельминтозы каб. 26', 'Анализ крови и мазки из зева и носа, на гельминтозы', 'analysis', NULL, true, false, false);
INSERT INTO public.doctor (id, title, description, doctor_type, sex, is_active, is_fun_male, is_fun_female) VALUES (12, 'ЭКГ каб. 7', 'Электрокардиография', 'analysis', NULL, true, false, false);
INSERT INTO public.doctor (id, title, description, doctor_type, sex, is_active, is_fun_male, is_fun_female) VALUES (13, 'Флюорография каб. 43', 'Флюорография в 2-х проекциях легких', 'analysis', NULL, true, false, false);
INSERT INTO public.doctor (id, title, description, doctor_type, sex, is_active, is_fun_male, is_fun_female) VALUES (9, 'Терапевт каб. 5 или 8', 'Врач - Терапевт', 'therapist', NULL, true, false, false);
INSERT INTO public.doctor (id, title, description, doctor_type, sex, is_active, is_fun_male, is_fun_female) VALUES (23, 'Анализ кала и мочи каб. 222 эт. 2', 'Анализ кала и мочи', 'analysis', NULL, true, false, false);


--
-- Data for Name: sticker; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.sticker (id, name, file, is_active) VALUES (3, 'hello', 'CAACAgIAAxUAAWJ6uOWBCOSuZ3LN-uyVewuR62WxAAI2FgACcmugS6XaTV2HP2QpJAQ', true);
INSERT INTO public.sticker (id, name, file, is_active) VALUES (4, 'omg', 'CAACAgIAAxUAAWJ6uOX-kZcyAAEMwLBrYODaf_C4ugACKhcAAktUoEuC6yy0FAABmRkkBA', true);
INSERT INTO public.sticker (id, name, file, is_active) VALUES (5, 'super', 'CAACAgIAAxUAAWJ6uOVJeBzxB_Ot1VI4woV9zS-CAAISFQACXis5SEb-mcPIFV6SJAQ', true);


--
-- Data for Name: telegram; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Data for Name: user; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Data for Name: user_doctor_done; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Name: newtable_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.newtable_id_seq', 24, true);


--
-- Name: stiker_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.stiker_id_seq', 5, true);


--
-- Name: telegram_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.telegram_id_seq', 1, true);


--
-- Name: user_doctor_done_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.user_doctor_done_id_seq', 692, true);


--
-- Name: user_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.user_id_seq', 70, true);


--
-- Data for Name: BLOBS; Type: BLOBS; Schema: -; Owner: 
--



--
-- Name: doctor doctor_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.doctor
    ADD CONSTRAINT doctor_pk PRIMARY KEY (id);


--
-- Name: sticker stiker_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sticker
    ADD CONSTRAINT stiker_pk PRIMARY KEY (id);


--
-- Name: user_doctor_done user_doctor_done_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_doctor_done
    ADD CONSTRAINT user_doctor_done_pk PRIMARY KEY (id);


--
-- Name: user user_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."user"
    ADD CONSTRAINT user_pk PRIMARY KEY (id);


--
-- Name: user_therapist_show _RETURN; Type: RULE; Schema: public; Owner: postgres
--

CREATE OR REPLACE VIEW public.user_therapist_show AS
 WITH RECURSIVE user_therapist_show(id, chat_id, title, show, is_active) AS (
         SELECT d_out.id,
            res.chat_id,
            d_out.title,
            bool_and(res.ready) AS show,
            d_out.is_active
           FROM public.doctor d_out,
            ( SELECT user_analysis_to_go.id,
                    user_analysis_to_go.chat_id,
                    user_analysis_to_go.title,
                    user_analysis_to_go.ready
                   FROM public.user_analysis_to_go
                UNION
                 SELECT user_doctor_to_go.id,
                    user_doctor_to_go.chat_id,
                    user_doctor_to_go.title,
                    user_doctor_to_go.ready
                   FROM public.user_doctor_to_go) res
          WHERE ((d_out.doctor_type = 'therapist'::public.doctor_type) AND d_out.is_active)
          GROUP BY d_out.id, res.chat_id
        )
 SELECT user_therapist_show.id,
    user_therapist_show.chat_id,
    user_therapist_show.title,
    user_therapist_show.show,
    user_therapist_show.is_active
   FROM user_therapist_show;


--
-- Name: user_doctor_done user_doctor_done_doctor_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_doctor_done
    ADD CONSTRAINT user_doctor_done_doctor_id_fkey FOREIGN KEY (doctor_id) REFERENCES public.doctor(id);


--
-- Name: user_doctor_done user_doctor_done_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_doctor_done
    ADD CONSTRAINT user_doctor_done_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id);


--
-- PostgreSQL database dump complete
--

