CREATE TYPE public.user_role AS ENUM('read', 'write', 'owner');

CREATE TABLE IF NOT EXISTS public.project_members (
	id bigserial NOT NULL,
	user_id integer NOT NULL,
	project_id integer NOT NULL,
	environment_id integer NULL,
	role public.user_role NOT NULL,
	project_owner boolean DEFAULT FALSE,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	CONSTRAINT project_members_pkey PRIMARY KEY (id),
	UNIQUE (user_id, project_id, environment_id)
);

DROP INDEX IF EXISTS idx_project_members_project_id;
CREATE INDEX idx_project_members_project_id ON public.project_members USING btree (project_id);
DROP INDEX IF EXISTS idx_project_members_user_id;
CREATE INDEX idx_project_members_user_id ON public.project_members USING btree (user_id);
DROP INDEX IF EXISTS idx_project_members_environment_id;
CREATE INDEX idx_project_members_environment_id ON public.project_members USING btree (environment_id);

-- public.environment_user_secrets foreign keys

ALTER TABLE public.project_members ADD CONSTRAINT fk_project_members_project FOREIGN KEY (project_id) REFERENCES projects(id);
ALTER TABLE public.project_members ADD CONSTRAINT fk_project_members_user FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE public.project_members ADD CONSTRAINT fk_project_members_environment FOREIGN KEY (environment_id) REFERENCES environments(id);
