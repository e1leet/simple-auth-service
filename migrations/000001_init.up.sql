BEGIN;

SET statement_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = ON;
SET check_function_bodies = FALSE;
SET client_min_messages = WARNING;
SET search_path = public, extensions;
SET default_tablespace = '';
SET default_with_oids = FALSE;

-- TABLES --
CREATE TABLE public.usr
(
    id         serial PRIMARY KEY,
    username   varchar(32) UNIQUE       NOT NULL,
    password   varchar(64)              NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now()
);

COMMIT;