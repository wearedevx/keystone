DROP INDEX IF EXISTS idx_environments_project_id;
ALTER TABLE public.environments DROP CONSTRAINT fk_projects_environments;

DROP TABLE IF EXISTS environments;
