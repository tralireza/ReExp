\connect reexp

--
-- PostgreSQL database dump
--

-- Dumped from database version 16.3 (Debian 16.3-1.pgdg120+1)
-- Dumped by pg_dump version 16.3 (Debian 16.3-1.pgdg120+1)

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

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: cfp; Type: TABLE; Schema: public; Owner: reexp
--

CREATE TABLE public.cfp (
    clinet integer NOT NULL,
    fund integer NOT NULL,
    portfolio integer NOT NULL
);


ALTER TABLE public.cfp OWNER TO reexp;

--
-- Name: client; Type: TABLE; Schema: public; Owner: reexp
--

CREATE TABLE public.client (
    id integer NOT NULL,
    dob date NOT NULL,
    name character varying(255) NOT NULL,
    ni character(9) NOT NULL,
    CONSTRAINT client_name_check CHECK (((name)::text <> ''::text)),
    CONSTRAINT client_ni_check CHECK (((ni <> ''::bpchar) AND (length(ni) = 9)))
);


ALTER TABLE public.client OWNER TO reexp;

--
-- Name: client_id_seq; Type: SEQUENCE; Schema: public; Owner: reexp
--

CREATE SEQUENCE public.client_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.client_id_seq OWNER TO reexp;

--
-- Name: client_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: reexp
--

ALTER SEQUENCE public.client_id_seq OWNED BY public.client.id;


--
-- Name: fund; Type: TABLE; Schema: public; Owner: reexp
--

CREATE TABLE public.fund (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    sector character varying(255) NOT NULL,
    type character varying(255) NOT NULL,
    CONSTRAINT fund_name_check CHECK (((name)::text <> ''::text)),
    CONSTRAINT fund_sector_check CHECK (((sector)::text <> ''::text)),
    CONSTRAINT fund_type_check CHECK (((type)::text <> ''::text))
);


ALTER TABLE public.fund OWNER TO reexp;

--
-- Name: fund_id_seq; Type: SEQUENCE; Schema: public; Owner: reexp
--

CREATE SEQUENCE public.fund_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.fund_id_seq OWNER TO reexp;

--
-- Name: fund_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: reexp
--

ALTER SEQUENCE public.fund_id_seq OWNED BY public.fund.id;


--
-- Name: portfolio; Type: TABLE; Schema: public; Owner: reexp
--

CREATE TABLE public.portfolio (
    id integer NOT NULL,
    amount numeric(6,0) NOT NULL,
    name character varying(255) NOT NULL,
    openedat date DEFAULT now() NOT NULL,
    state integer DEFAULT 0 NOT NULL,
    CONSTRAINT portfolio_amount_check CHECK ((amount > (0)::numeric)),
    CONSTRAINT portfolio_name_check CHECK (((name)::text <> ''::text))
);


ALTER TABLE public.portfolio OWNER TO reexp;

--
-- Name: portfolio_id_seq; Type: SEQUENCE; Schema: public; Owner: reexp
--

CREATE SEQUENCE public.portfolio_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.portfolio_id_seq OWNER TO reexp;

--
-- Name: portfolio_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: reexp
--

ALTER SEQUENCE public.portfolio_id_seq OWNED BY public.portfolio.id;


--
-- Name: client id; Type: DEFAULT; Schema: public; Owner: reexp
--

ALTER TABLE ONLY public.client ALTER COLUMN id SET DEFAULT nextval('public.client_id_seq'::regclass);


--
-- Name: fund id; Type: DEFAULT; Schema: public; Owner: reexp
--

ALTER TABLE ONLY public.fund ALTER COLUMN id SET DEFAULT nextval('public.fund_id_seq'::regclass);


--
-- Name: portfolio id; Type: DEFAULT; Schema: public; Owner: reexp
--

ALTER TABLE ONLY public.portfolio ALTER COLUMN id SET DEFAULT nextval('public.portfolio_id_seq'::regclass);


--
-- Name: cfp cfp_pkey; Type: CONSTRAINT; Schema: public; Owner: reexp
--

ALTER TABLE ONLY public.cfp
    ADD CONSTRAINT cfp_pkey PRIMARY KEY (clinet);


--
-- Name: client client_ni_key; Type: CONSTRAINT; Schema: public; Owner: reexp
--

ALTER TABLE ONLY public.client
    ADD CONSTRAINT client_ni_key UNIQUE (ni);


--
-- Name: client client_pkey; Type: CONSTRAINT; Schema: public; Owner: reexp
--

ALTER TABLE ONLY public.client
    ADD CONSTRAINT client_pkey PRIMARY KEY (id);


--
-- Name: fund fund_pkey; Type: CONSTRAINT; Schema: public; Owner: reexp
--

ALTER TABLE ONLY public.fund
    ADD CONSTRAINT fund_pkey PRIMARY KEY (id);


--
-- Name: portfolio portfolio_pkey; Type: CONSTRAINT; Schema: public; Owner: reexp
--

ALTER TABLE ONLY public.portfolio
    ADD CONSTRAINT portfolio_pkey PRIMARY KEY (id);


--
-- Name: cfp cfp_fund_fkey; Type: FK CONSTRAINT; Schema: public; Owner: reexp
--

ALTER TABLE ONLY public.cfp
    ADD CONSTRAINT cfp_fund_fkey FOREIGN KEY (fund) REFERENCES public.fund(id);


--
-- Name: cfp cfp_portfolio_fkey; Type: FK CONSTRAINT; Schema: public; Owner: reexp
--

ALTER TABLE ONLY public.cfp
    ADD CONSTRAINT cfp_portfolio_fkey FOREIGN KEY (portfolio) REFERENCES public.portfolio(id);


--
-- PostgreSQL database dump complete
--

