ALTER TABLE public.projects
DROP CONSTRAINT IF EXISTS fk_projects_user;

DROP INDEX IF EXISTS idx_projects_user_id;

ALTER TABLE public.projects
DROP COLUMN IF EXISTS user_id;

-- 
ALTER TABLE public.environments
DROP CONSTRAINT IF EXISTS fk_environments_project;

DROP INDEX IF EXISTS idx_environments_project_id;

ALTER TABLE public.environments
DROP COLUMN IF EXISTS project_id,
DROP COLUMN IF EXISTS versionID;

--
ALTER TABLE public.project_members
DROP CONSTRAINT IF EXISTS fk_project_members_role;

DROP INDEX IF EXISTS idx_project_members_role_id;

ALTER TABLE public.project_members
DROP COLUMN IF EXISTS role_id,
ADD COLUMN IF NOT EXISTS role public.user_role NOT NULL,
ADD COLUMN IF NOT EXISTS project_owner bool NOT NULL default false,
ADD COLUMN IF NOT EXISTS environment_id integer;

CREATE INDEX IF NOT EXISTS idx_project_members_environment_id ON public.project_members USING btree (environment_id);

DELETE FROM public.project_members;

ALTER TABLE public.project_members
DROP CONSTRAINT IF EXISTS project_members_user_id_project_id_key,
DROP CONSTRAINT IF EXISTS project_members_user_id_project_id_environment_id_key,
DROP CONSTRAINT IF EXISTS fk_project_members_environment,
ADD CONSTRAINT project_members_user_id_project_id_environment_id_key UNIQUE (user_id, project_id, environment_id),
ADD CONSTRAINT fk_project_members_environment FOREIGN KEY (environment_id) REFERENCES public.environments(id);

-- Environments
ALTER TABLE public.environments
DROP CONSTRAINT IF EXISTS fk_environments_environment_type;

DROP INDEX IF EXISTS idx_environments_environment_type_id;

ALTER TABLE public.environments
DROP COLUMN IF EXISTS environment_type_id;

-- Relation between environment_types and roles

DROP INDEX IF EXISTS idx_role_environment_types_role_id;
DROP INDEX IF EXISTS idx_role_environment_types_environment_type_id;

ALTER TABLE IF EXISTS public.roles_environment_types
DROP CONSTRAINT IF EXISTS fk_role_environment_types_role,
DROP CONSTRAINT IF EXISTS fk_role_environment_types_environment_type;

DROP TABLE IF EXISTS public.roles_environment_types;

-- Environment type
DROP TABLE IF EXISTS public.environment_types;

-- Roles
DROP TABLE IF EXISTS public.roles;
