DROP INDEX IF EXISTS idx_environments_project_id;
ALTER TABLE public.environments DROP CONSTRAINT fk_projects_environments;
ALTER TABLE public.environments DROP COLUMN project_id;

ALTER TABLE public.users DROP COLUMN keys_cipher; 
ALTER TABLE public.users DROP COLUMN keys_sign;
