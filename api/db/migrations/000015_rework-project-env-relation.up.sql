-- Role table

CREATE TABLE IF NOT EXISTS public.roles (
	id bigserial NOT NULL,
	name varchar(255) NOT NULL,
	description text,
	created_at timestamptz,
	updated_at timestamptz,
	CONSTRAINT roles_pkey PRIMARY KEY (id)
);

-- Environment types
CREATE TABLE IF NOT EXISTS public.environment_types (
	id bigserial NOT NULL,
	name varchar(255) NOT NULL,
	created_at timestamptz,
	updated_at timestamptz,
	CONSTRAINT environment_types_pkey PRIMARY KEY (id)
);

-- Relation between environment_types and roles
CREATE TABLE IF NOT EXISTS public.roles_environment_types (
	id bigserial NOT NULL,
	role_id integer NOT NULL,
	environment_type_id integer NOT NULL,
	read bool NOT NULL default false,
	write bool NOT NULL default false,
	invite bool NOT NULL default false,
	created_at timestamptz,
	updated_at timestamptz,
	CONSTRAINT roles_environment_types_pkey PRIMARY KEY (id)
);

DROP INDEX IF EXISTS idx_role_environment_types_role_id;
CREATE INDEX  idx_role_environment_types_role_id ON public.roles_environment_types USING btree (role_id);
DROP INDEX IF EXISTS idx_role_environment_types_environment_type_id;
CREATE INDEX  idx_role_environment_types_environment_type_id ON public.roles_environment_types USING btree (environment_type_id);

ALTER TABLE public.roles_environment_types
ADD CONSTRAINT fk_role_environment_types_role FOREIGN KEY (role_id) REFERENCES public.roles(id),
ADD CONSTRAINT fk_role_environment_types_environment_type FOREIGN KEY (environment_type_id) REFERENCES public.environment_types(id);

-- Environment

ALTER TABLE public.environments
ADD COLUMN environment_type_id integer;

DROP INDEX IF EXISTS idx_environments_environment_type_id;
CREATE INDEX idx_environments_environment_type_id ON public.environments USING btree (environment_type_id);

ALTER TABLE public.environments
ADD CONSTRAINT fk_environments_environment_type FOREIGN KEY (environment_type_id) REFERENCES public.environment_types(id);


-- Project Member
-- remove useless relationships

ALTER TABLE public.project_members
DROP CONSTRAINT fk_project_members_environment,
DROP CONSTRAINT project_members_user_id_project_id_environment_id_key; 

DROP INDEX IF EXISTS project_members_user_id_project_id_environment_id_key; 

DROP INDEX idx_project_members_environment_id;

ALTER TABLE public.project_members
DROP COLUMN environment_id,
DROP COLUMN role,
DROP COLUMN project_owner;

ALTER TABLE public.project_members
DROP CONSTRAINT IF EXISTS project_members_user_id_project_id_key;

DROP INDEX IF EXISTS project_members_user_id_project_id_key;

DELETE FROM public.project_members;

ALTER TABLE public.project_members
ADD CONSTRAINT project_members_user_id_project_id_key UNIQUE (user_id, project_id);

-- Relation between project members and roles
ALTER TABLE public.project_members
ADD COLUMN role_id integer;

DROP INDEX IF EXISTS idx_project_members_role_id;
CREATE INDEX idx_project_members_role_id ON public.project_members USING btree (role_id);

ALTER TABLE public.project_members
ADD CONSTRAINT fk_project_members_role FOREIGN KEY (role_id) REFERENCES public.roles(id);

-- Environment belongs to project
ALTER TABLE public.environments
ADD COLUMN IF NOT EXISTS version_id varchar(255),
ADD COLUMN IF NOT EXISTS project_id integer;

DROP INDEX IF EXISTS idx_environments_project_id;
CREATE INDEX idx_environments_project_id ON public.environments USING btree (project_id);

ALTER TABLE public.environments
ADD CONSTRAINT fk_environments_project FOREIGN KEY (project_id) REFERENCES projects(id);

-- Project belongs to one user
DELETE FROM public.projects;

ALTER TABLE public.projects
ADD COLUMN IF NOT EXISTS user_id integer NOT NULL;

DROP INDEX IF EXISTS idx_projects_user_id;
CREATE INDEX idx_projects_user_id ON public.projects USING btree (user_id);

ALTER TABLE public.projects
DROP CONSTRAINT IF EXISTS fk_projects_user,
ADD CONSTRAINT fk_projects_user FOREIGN KEY (user_id) REFERENCES public.users(id);
