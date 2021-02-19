DROP INDEX IF EXISTS idx_project_members_project_id;
DROP INDEX IF EXISTS idx_project_members_user_id;
DROP INDEX IF EXISTS idx_project_members_environment_id;

ALTER TABLE public.project_members DROP CONSTRAINT fk_project_members_environment;
ALTER TABLE public.project_members DROP CONSTRAINT fk_project_members_project;
ALTER TABLE public.project_members DROP CONSTRAINT fk_project_members_user;

DROP TABLE IF EXISTS public.project_members;

DROP TYPE IF EXISTS public.user_role;
