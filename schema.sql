--
-- PostgreSQL database dump
--

\restrict yPFSNSz2D9nQMQ9484kNoS40rEM5mu7dvbWKKUbrYQPKc4MKknXz5G22P067g7B

-- Dumped from database version 15.17
-- Dumped by pg_dump version 15.17

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
-- Name: public; Type: SCHEMA; Schema: -; Owner: pg_database_owner
--

CREATE SCHEMA public;


ALTER SCHEMA public OWNER TO pg_database_owner;

--
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: pg_database_owner
--

COMMENT ON SCHEMA public IS 'standard public schema';


--
-- Name: update_current_timestamp(); Type: FUNCTION; Schema: public; Owner: nanomdm
--

CREATE FUNCTION public.update_current_timestamp() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = now();
RETURN NEW;
END;
$$;


ALTER FUNCTION public.update_current_timestamp() OWNER TO nanomdm;

--
-- Name: update_updated_at(); Type: FUNCTION; Schema: public; Owner: nanomdm
--

CREATE FUNCTION public.update_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.update_updated_at() OWNER TO nanomdm;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: cert_auth_associations; Type: TABLE; Schema: public; Owner: nanomdm
--

CREATE TABLE public.cert_auth_associations (
    id character varying(255) NOT NULL,
    sha256 character(64) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT cert_auth_associations_id_check CHECK (((id)::text <> ''::text)),
    CONSTRAINT cert_auth_associations_sha256_check CHECK ((sha256 <> ''::bpchar))
);


ALTER TABLE public.cert_auth_associations OWNER TO nanomdm;

--
-- Name: command_results; Type: TABLE; Schema: public; Owner: nanomdm
--

CREATE TABLE public.command_results (
    id character varying(255) NOT NULL,
    command_uuid character varying(127) NOT NULL,
    status character varying(31) NOT NULL,
    result text NOT NULL,
    not_now_at timestamp without time zone,
    not_now_tally integer DEFAULT 0 NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT command_results_result_check CHECK ((SUBSTRING(result FROM 1 FOR 5) = '<?xml'::text)),
    CONSTRAINT command_results_status_check CHECK (((status)::text <> ''::text))
);


ALTER TABLE public.command_results OWNER TO nanomdm;

--
-- Name: commands; Type: TABLE; Schema: public; Owner: nanomdm
--

CREATE TABLE public.commands (
    command_uuid character varying(127) NOT NULL,
    request_type character varying(63) NOT NULL,
    command text NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT commands_command_check CHECK ((SUBSTRING(command FROM 1 FOR 5) = '<?xml'::text)),
    CONSTRAINT commands_command_uuid_check CHECK (((command_uuid)::text <> ''::text)),
    CONSTRAINT commands_request_type_check CHECK (((request_type)::text <> ''::text))
);


ALTER TABLE public.commands OWNER TO nanomdm;

--
-- Name: dep_names; Type: TABLE; Schema: public; Owner: nanomdm
--

CREATE TABLE public.dep_names (
    name character varying(255) NOT NULL,
    consumer_key text,
    consumer_secret text,
    access_token text,
    access_secret text,
    access_token_expiry timestamp with time zone,
    config_base_url character varying(255),
    tokenpki_cert_pem text,
    tokenpki_key_pem text,
    tokenpki_staging_cert_pem text,
    tokenpki_staging_key_pem text,
    syncer_cursor character varying(1024),
    assigner_profile_uuid text,
    assigner_profile_uuid_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT dep_names_tokenpki_cert_pem_check CHECK (((tokenpki_cert_pem IS NULL) OR (SUBSTRING(tokenpki_cert_pem FROM 1 FOR 27) = '-----BEGIN CERTIFICATE-----'::text))),
    CONSTRAINT dep_names_tokenpki_key_pem_check CHECK (((tokenpki_key_pem IS NULL) OR (SUBSTRING(tokenpki_key_pem FROM 1 FOR 5) = '-----'::text)))
);


ALTER TABLE public.dep_names OWNER TO nanomdm;

--
-- Name: devices; Type: TABLE; Schema: public; Owner: nanomdm
--

CREATE TABLE public.devices (
    id character varying(255) NOT NULL,
    identity_cert text,
    serial_number character varying(127),
    unlock_token bytea,
    unlock_token_at timestamp without time zone,
    authenticate text NOT NULL,
    authenticate_at timestamp without time zone NOT NULL,
    token_update text,
    token_update_at timestamp without time zone,
    bootstrap_token_b64 text,
    bootstrap_token_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT devices_authenticate_check CHECK ((authenticate <> ''::text)),
    CONSTRAINT devices_bootstrap_token_b64_check CHECK (((bootstrap_token_b64 IS NULL) OR (bootstrap_token_b64 <> ''::text))),
    CONSTRAINT devices_identity_cert_check CHECK (((identity_cert IS NULL) OR (SUBSTRING(identity_cert FROM 1 FOR 27) = '-----BEGIN CERTIFICATE-----'::text))),
    CONSTRAINT devices_serial_number_check CHECK (((serial_number IS NULL) OR ((serial_number)::text <> ''::text))),
    CONSTRAINT devices_token_update_check CHECK (((token_update IS NULL) OR (token_update <> ''::text))),
    CONSTRAINT devices_unlock_token_check CHECK (((unlock_token IS NULL) OR (length(unlock_token) > 0)))
);


ALTER TABLE public.devices OWNER TO nanomdm;

--
-- Name: enrollment_queue; Type: TABLE; Schema: public; Owner: nanomdm
--

CREATE TABLE public.enrollment_queue (
    id character varying(255) NOT NULL,
    command_uuid character varying(127) NOT NULL,
    active boolean DEFAULT true NOT NULL,
    priority smallint DEFAULT 0 NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.enrollment_queue OWNER TO nanomdm;

--
-- Name: enrollments; Type: TABLE; Schema: public; Owner: nanomdm
--

CREATE TABLE public.enrollments (
    id character varying(255) NOT NULL,
    device_id character varying(255) NOT NULL,
    user_id character varying(255),
    type character varying(31) NOT NULL,
    topic character varying(255) NOT NULL,
    push_magic character varying(127) NOT NULL,
    token_hex character varying(255) NOT NULL,
    enabled boolean DEFAULT true NOT NULL,
    token_update_tally integer DEFAULT 1 NOT NULL,
    last_seen_at timestamp without time zone NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT enrollments_id_check CHECK (((id)::text <> ''::text)),
    CONSTRAINT enrollments_push_magic_check CHECK (((push_magic)::text <> ''::text)),
    CONSTRAINT enrollments_token_hex_check CHECK (((token_hex)::text <> ''::text)),
    CONSTRAINT enrollments_topic_check CHECK (((topic)::text <> ''::text)),
    CONSTRAINT enrollments_type_check CHECK (((type)::text <> ''::text))
);


ALTER TABLE public.enrollments OWNER TO nanomdm;

--
-- Name: push_certs; Type: TABLE; Schema: public; Owner: nanomdm
--

CREATE TABLE public.push_certs (
    topic character varying(255) NOT NULL,
    cert_pem text NOT NULL,
    key_pem text NOT NULL,
    stale_token integer NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT push_certs_cert_pem_check CHECK ((SUBSTRING(cert_pem FROM 1 FOR 27) = '-----BEGIN CERTIFICATE-----'::text)),
    CONSTRAINT push_certs_key_pem_check CHECK ((SUBSTRING(key_pem FROM 1 FOR 5) = '-----'::text)),
    CONSTRAINT push_certs_topic_check CHECK (((topic)::text <> ''::text))
);


ALTER TABLE public.push_certs OWNER TO nanomdm;

--
-- Name: users; Type: TABLE; Schema: public; Owner: nanomdm
--

CREATE TABLE public.users (
    id character varying(255) NOT NULL,
    device_id character varying(255) NOT NULL,
    user_short_name character varying(255),
    user_long_name character varying(255),
    token_update text,
    token_update_at timestamp without time zone,
    user_authenticate text,
    user_authenticate_at timestamp without time zone,
    user_authenticate_digest text,
    user_authenticate_digest_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT users_token_update_check CHECK (((token_update IS NULL) OR (token_update <> ''::text))),
    CONSTRAINT users_user_authenticate_check CHECK (((user_authenticate IS NULL) OR (user_authenticate <> ''::text))),
    CONSTRAINT users_user_authenticate_digest_check CHECK (((user_authenticate_digest IS NULL) OR (user_authenticate_digest <> ''::text))),
    CONSTRAINT users_user_long_name_check CHECK (((user_long_name IS NULL) OR ((user_long_name)::text <> ''::text))),
    CONSTRAINT users_user_short_name_check CHECK (((user_short_name IS NULL) OR ((user_short_name)::text <> ''::text)))
);


ALTER TABLE public.users OWNER TO nanomdm;

--
-- Name: view_queue; Type: VIEW; Schema: public; Owner: nanomdm
--

CREATE VIEW public.view_queue AS
 SELECT q.id,
    q.created_at,
    q.active,
    q.priority,
    c.command_uuid,
    c.request_type,
    c.command,
    r.updated_at AS result_updated_at,
    r.status,
    r.result
   FROM ((public.enrollment_queue q
     JOIN public.commands c ON (((q.command_uuid)::text = (c.command_uuid)::text)))
     LEFT JOIN public.command_results r ON ((((r.command_uuid)::text = (q.command_uuid)::text) AND ((r.id)::text = (q.id)::text))))
  ORDER BY q.priority DESC, q.created_at;


ALTER TABLE public.view_queue OWNER TO nanomdm;

--
-- Name: cert_auth_associations cert_auth_associations_pkey; Type: CONSTRAINT; Schema: public; Owner: nanomdm
--

ALTER TABLE ONLY public.cert_auth_associations
    ADD CONSTRAINT cert_auth_associations_pkey PRIMARY KEY (id, sha256);


--
-- Name: command_results command_results_pkey; Type: CONSTRAINT; Schema: public; Owner: nanomdm
--

ALTER TABLE ONLY public.command_results
    ADD CONSTRAINT command_results_pkey PRIMARY KEY (id, command_uuid);


--
-- Name: commands commands_pkey; Type: CONSTRAINT; Schema: public; Owner: nanomdm
--

ALTER TABLE ONLY public.commands
    ADD CONSTRAINT commands_pkey PRIMARY KEY (command_uuid);


--
-- Name: dep_names dep_names_pkey; Type: CONSTRAINT; Schema: public; Owner: nanomdm
--

ALTER TABLE ONLY public.dep_names
    ADD CONSTRAINT dep_names_pkey PRIMARY KEY (name);


--
-- Name: devices devices_pkey; Type: CONSTRAINT; Schema: public; Owner: nanomdm
--

ALTER TABLE ONLY public.devices
    ADD CONSTRAINT devices_pkey PRIMARY KEY (id);


--
-- Name: enrollment_queue enrollment_queue_pkey; Type: CONSTRAINT; Schema: public; Owner: nanomdm
--

ALTER TABLE ONLY public.enrollment_queue
    ADD CONSTRAINT enrollment_queue_pkey PRIMARY KEY (id, command_uuid);


--
-- Name: enrollments enrollments_pkey; Type: CONSTRAINT; Schema: public; Owner: nanomdm
--

ALTER TABLE ONLY public.enrollments
    ADD CONSTRAINT enrollments_pkey PRIMARY KEY (id);


--
-- Name: enrollments enrollments_user_id_key; Type: CONSTRAINT; Schema: public; Owner: nanomdm
--

ALTER TABLE ONLY public.enrollments
    ADD CONSTRAINT enrollments_user_id_key UNIQUE (user_id);


--
-- Name: push_certs push_certs_pkey; Type: CONSTRAINT; Schema: public; Owner: nanomdm
--

ALTER TABLE ONLY public.push_certs
    ADD CONSTRAINT push_certs_pkey PRIMARY KEY (topic);


--
-- Name: users users_id_key; Type: CONSTRAINT; Schema: public; Owner: nanomdm
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_id_key UNIQUE (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: nanomdm
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id, device_id);


--
-- Name: idx_status; Type: INDEX; Schema: public; Owner: nanomdm
--

CREATE INDEX idx_status ON public.command_results USING btree (status);


--
-- Name: idx_type; Type: INDEX; Schema: public; Owner: nanomdm
--

CREATE INDEX idx_type ON public.enrollments USING btree (type);


--
-- Name: serial_number; Type: INDEX; Schema: public; Owner: nanomdm
--

CREATE INDEX serial_number ON public.devices USING btree (serial_number);


--
-- Name: cert_auth_associations update_at_to_current_timestamp; Type: TRIGGER; Schema: public; Owner: nanomdm
--

CREATE TRIGGER update_at_to_current_timestamp BEFORE UPDATE ON public.cert_auth_associations FOR EACH ROW EXECUTE FUNCTION public.update_current_timestamp();


--
-- Name: command_results update_at_to_current_timestamp; Type: TRIGGER; Schema: public; Owner: nanomdm
--

CREATE TRIGGER update_at_to_current_timestamp BEFORE UPDATE ON public.command_results FOR EACH ROW EXECUTE FUNCTION public.update_current_timestamp();


--
-- Name: commands update_at_to_current_timestamp; Type: TRIGGER; Schema: public; Owner: nanomdm
--

CREATE TRIGGER update_at_to_current_timestamp BEFORE UPDATE ON public.commands FOR EACH ROW EXECUTE FUNCTION public.update_current_timestamp();


--
-- Name: devices update_at_to_current_timestamp; Type: TRIGGER; Schema: public; Owner: nanomdm
--

CREATE TRIGGER update_at_to_current_timestamp BEFORE UPDATE ON public.devices FOR EACH ROW EXECUTE FUNCTION public.update_current_timestamp();


--
-- Name: enrollment_queue update_at_to_current_timestamp; Type: TRIGGER; Schema: public; Owner: nanomdm
--

CREATE TRIGGER update_at_to_current_timestamp BEFORE UPDATE ON public.enrollment_queue FOR EACH ROW EXECUTE FUNCTION public.update_current_timestamp();


--
-- Name: enrollments update_at_to_current_timestamp; Type: TRIGGER; Schema: public; Owner: nanomdm
--

CREATE TRIGGER update_at_to_current_timestamp BEFORE UPDATE ON public.enrollments FOR EACH ROW EXECUTE FUNCTION public.update_current_timestamp();


--
-- Name: push_certs update_at_to_current_timestamp; Type: TRIGGER; Schema: public; Owner: nanomdm
--

CREATE TRIGGER update_at_to_current_timestamp BEFORE UPDATE ON public.push_certs FOR EACH ROW EXECUTE FUNCTION public.update_current_timestamp();


--
-- Name: users update_at_to_current_timestamp; Type: TRIGGER; Schema: public; Owner: nanomdm
--

CREATE TRIGGER update_at_to_current_timestamp BEFORE UPDATE ON public.users FOR EACH ROW EXECUTE FUNCTION public.update_current_timestamp();


--
-- Name: dep_names update_updated_at_on_change; Type: TRIGGER; Schema: public; Owner: nanomdm
--

CREATE TRIGGER update_updated_at_on_change BEFORE UPDATE ON public.dep_names FOR EACH ROW EXECUTE FUNCTION public.update_updated_at();


--
-- Name: command_results command_results_command_uuid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: nanomdm
--

ALTER TABLE ONLY public.command_results
    ADD CONSTRAINT command_results_command_uuid_fkey FOREIGN KEY (command_uuid) REFERENCES public.commands(command_uuid) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: command_results command_results_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: nanomdm
--

ALTER TABLE ONLY public.command_results
    ADD CONSTRAINT command_results_id_fkey FOREIGN KEY (id) REFERENCES public.enrollments(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: enrollment_queue enrollment_queue_command_uuid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: nanomdm
--

ALTER TABLE ONLY public.enrollment_queue
    ADD CONSTRAINT enrollment_queue_command_uuid_fkey FOREIGN KEY (command_uuid) REFERENCES public.commands(command_uuid) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: enrollment_queue enrollment_queue_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: nanomdm
--

ALTER TABLE ONLY public.enrollment_queue
    ADD CONSTRAINT enrollment_queue_id_fkey FOREIGN KEY (id) REFERENCES public.enrollments(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: enrollments enrollments_device_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: nanomdm
--

ALTER TABLE ONLY public.enrollments
    ADD CONSTRAINT enrollments_device_id_fkey FOREIGN KEY (device_id) REFERENCES public.devices(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: enrollments enrollments_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: nanomdm
--

ALTER TABLE ONLY public.enrollments
    ADD CONSTRAINT enrollments_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: users users_device_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: nanomdm
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_device_id_fkey FOREIGN KEY (device_id) REFERENCES public.devices(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

\unrestrict yPFSNSz2D9nQMQ9484kNoS40rEM5mu7dvbWKKUbrYQPKc4MKknXz5G22P067g7B

