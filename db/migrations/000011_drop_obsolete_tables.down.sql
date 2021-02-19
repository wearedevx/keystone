-- public.project_permissions

CREATE TABLE public.project_permissions (
	user_id integer NOT NULL,
  project_id integer NOT NULL,
	"role" text NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	CONSTRAINT project_permissions_pkey PRIMARY KEY (user_id, project_id)
);
-- public.project_permissions foreign keys
ALTER TABLE public.project_permissions ADD CONSTRAINT fk_project_permissions_project FOREIGN KEY (project_id) REFERENCES projects(id);
ALTER TABLE public.project_permissions ADD CONSTRAINT fk_project_permissions_user FOREIGN KEY (user_id) REFERENCES users(id);

-- public.environment_permissions
CREATE TABLE public.environment_permissions (
	user_id integer NOT NULL,
	environment_id integer NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	"role" text NULL,
	CONSTRAINT environment_permissions_pkey PRIMARY KEY (user_id, environment_id)
);
-- public.environment_permissions foreign keys
ALTER TABLE public.environment_permissions ADD CONSTRAINT fk_environment_permissions_environment FOREIGN KEY (environment_id) REFERENCES environments(id);
ALTER TABLE public.environment_permissions ADD CONSTRAINT fk_environment_permissions_user FOREIGN KEY (user_id) REFERENCES users(id);

-- public.secrets
CREATE TABLE IF NOT EXISTS public.secrets (
	id bigserial NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	"name" text NOT NULL,
	"type" text NOT NULL,
	CONSTRAINT secrets_pkey PRIMARY KEY (id)
);

-- public.environment_user_secrets
CREATE TABLE IF NOT EXISTS public.environment_user_secrets (
	user_id integer NOT NULL,
	secret_id integer NOT NULL,
	environment_id integer NULL,
	value bytea NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	CONSTRAINT environment_user_secrets_pkey PRIMARY KEY (user_id, secret_id, environment_id)
);
-- public.environment_user_secrets foreign keys
ALTER TABLE public.environment_user_secrets ADD CONSTRAINT fk_environment_user_secrets_secret FOREIGN KEY (secret_id) REFERENCES secrets(id);
ALTER TABLE public.environment_user_secrets ADD CONSTRAINT fk_environment_user_secrets_user FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE public.environment_user_secrets ADD CONSTRAINT fk_environment_user_secrets_environment FOREIGN KEY (environment_id) REFERENCES environments(id);
