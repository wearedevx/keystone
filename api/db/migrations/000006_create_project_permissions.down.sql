ALTER TABLE public.project_permissions DROP CONSTRAINT fk_project_permissions_project;
ALTER TABLE public.project_permissions DROP CONSTRAINT fk_project_permissions_user;

DROP TABLE IF EXISTS public.project_permissions;
