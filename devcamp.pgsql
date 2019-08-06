--
-- PostgreSQL database dump
--

-- Dumped from database version 11.4
-- Dumped by pg_dump version 11.4

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

SET default_with_oids = false;

--
-- Name: books; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.books (
    id integer NOT NULL,
    title character varying,
    author character varying,
    isbn character varying,
    stock integer
);


ALTER TABLE public.books OWNER TO postgres;

--
-- Name: lend; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.lend (
    user_id integer,
    book_id integer
);


ALTER TABLE public.lend OWNER TO postgres;

--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id integer NOT NULL,
    name character varying
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Data for Name: books; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.books (id, title, author, isbn, stock) FROM stdin;
1	"The Adventures of Duck and Goose"	"Sir Quackalot"	"ISBN-TEST"	10
2	"The Return of Duck and Goose"	"Sir Quackalot"	"ISBN-TEST"	10
3	"More Fun with Duck and Goose"	"Sir Quackalot"	"ISBN-TEST"	10
4	"Duck and Goose on Holiday"	"Sir Quackalot"	"ISBN-TEST"	10
5	"The Return of Duck and Goose"	"Sir Quackalot"	"ISBN-TEST"	10
6	"The Adventures of Duck and Goose"	"Sir Quackalot"	"ISBN-TEST"	10
7	"My Friend is a Duck"	"A. Parrot"	"ISBN-TEST"	10
8	"Annotated Notes on the ‘Duck and Goose’ chronicles"	"Prof Macaw"	"ISBN-TEST"	10
9	"‘Duck and Goose’ Cheat Sheet for Students"	"Polly Parrot"	"ISBN-TEST"	10
10	"‘Duck and Goose’: an allegory for modern times?"	"Bor Ing"	"ISBN-TEST"	10
97	"‘Duck and Goose’: an allegory for modern times?"	\N	"ISBN-TEST"	10
\.


--
-- Data for Name: lend; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.lend (user_id, book_id) FROM stdin;
1	10
2	8
2	7
3	1
3	5
3	9
3	10
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, name) FROM stdin;
1	User 1
2	User 2
3	User Blabla
97	\N
98	\N
99	\N
\.


--
-- Name: books books_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.books
    ADD CONSTRAINT books_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: lend lend_book_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.lend
    ADD CONSTRAINT lend_book_id_fkey FOREIGN KEY (book_id) REFERENCES public.books(id);


--
-- Name: lend lend_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.lend
    ADD CONSTRAINT lend_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- PostgreSQL database dump complete
--

