-- drop  table table_name;
-- ALTER TABLE table_name DROP CONSTRAINT constraint_name; 
-- DROP TRIGGER trigger_name on table_name; 


-- uuid_generate_v4()
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';

ALTER DATABASE voting SET TIMEZONE TO 'America/Sao_Paulo';

SET client_encoding = 'UTF8';


CREATE TABLE public.t_wall (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz NULL,
    deleted_at timestamptz NULL,
    start_time timestamptz NULL,
    end_time timestamptz NULL,

);

COMMENT ON TABLE public.t_wall IS 'Esta tabela armazena os dados de um paredão.';

-- Permissions
GRANT INSERT, SELECT, UPDATE, DELETE ON TABLE public.t_wall TO voting_app;

ALTER TABLE ONLY public.t_wall
    ADD CONSTRAINT t_wall_pkey PRIMARY KEY (id);
    
CREATE TRIGGER t_wall_set_updated_at BEFORE UPDATE ON public.t_wall FOR EACH ROW EXECUTE PROCEDURE public.tf_utils_set_updated_at();


CREATE TABLE public.t_participant (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz NULL,
    deleted_at timestamptz NULL,
    name_participant text NOT NULL check(char_length(name_participant) <= 64)
);

COMMENT ON TABLE public.t_participant IS 'Esta tabela armazena os dados de um participante do paredão.';

-- Permissions
GRANT INSERT, SELECT, UPDATE, DELETE ON TABLE public.t_participant TO voting_app;

ALTER TABLE ONLY public.t_participant
    ADD CONSTRAINT t_participant_pkey PRIMARY KEY (id);

CREATE TRIGGER t_participant_set_updated_at BEFORE UPDATE ON public.t_participant FOR EACH ROW EXECUTE PROCEDURE public.tf_utils_set_updated_at();


CREATE TABLE public.t_wall_participant (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz NULL,
    deleted_at timestamptz NULL,
    wall_id uuid NOT NULL,
    participant_id uuid NOT NULL
);

COMMENT ON TABLE public.t_wall_participant IS 'Esta tabela armazena a relação dos participantes com o paredão.';

-- Permissions
GRANT INSERT, SELECT, UPDATE, DELETE ON TABLE public.t_wall_participant TO voting_app;

ALTER TABLE ONLY public.t_wall_participant
    ADD CONSTRAINT t_wall_participant_pkey PRIMARY KEY (id);
    
ALTER TABLE ONLY public.t_wall_participant
    ADD CONSTRAINT t_wall_participant_wall_id_fkey FOREIGN KEY (wall_id) REFERENCES public.t_wall(id);
    
ALTER TABLE ONLY public.t_wall_participant
    ADD CONSTRAINT t_wall_participant_participant_id_fkey FOREIGN KEY (participant_id) REFERENCES public.t_participant(id);

CREATE TRIGGER t_wall_participant_set_updated_at BEFORE UPDATE ON public.t_wall_participant FOR EACH ROW EXECUTE PROCEDURE public.tf_utils_set_updated_at();


CREATE TABLE public.t_vote (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz NULL,
    deleted_at timestamptz NULL,
    wall_id uuid NOT NULL,
    participant_id uuid NOT NULL,
    register_id uuid NOT NULL,
    register_at timestamptz NOT NULL,
    ip text NOT NULL
);

COMMENT ON TABLE public.t_vote IS 'Esta tabela armazena os votos de um paredão.';

-- Permissions
GRANT INSERT, SELECT, UPDATE, DELETE ON TABLE public.t_wall_participant TO voting_app;

ALTER TABLE ONLY public.t_vote
    ADD CONSTRAINT t_vote_pkey PRIMARY KEY (id);
    
ALTER TABLE ONLY public.t_vote
    ADD CONSTRAINT t_vote_wall_id_fkey FOREIGN KEY (wall_id) REFERENCES public.t_wall(id);
    
ALTER TABLE ONLY public.t_vote
    ADD CONSTRAINT t_vote_participant_id_fkey FOREIGN KEY (participant_id) REFERENCES public.t_participant(id);

CREATE TRIGGER t_vote_set_updated_at BEFORE UPDATE ON public.t_vote FOR EACH ROW EXECUTE PROCEDURE public.tf_utils_set_updated_at();

ALTER TABLE ONLY public.t_vote
    ADD CONSTRAINT register_id_unique UNIQUE (register_id);

