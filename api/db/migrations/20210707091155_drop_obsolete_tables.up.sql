-- project permission
ALTER TABLE public.project_permissions DROP CONSTRAINT fk_project_permissions_project;
ALTER TABLE public.project_permissions DROP CONSTRAINT fk_project_permissions_user;
DROP TABLE IF EXISTS public.project_permissions;

-- environment_permissionsons
ALTER TABLE public.environment_permissions DROP CONSTRAINT fk_environment_permissions_user;
ALTER TABLE public.environment_permissions DROP CONSTRAINT fk_environment_permissions_environment;
DROP TABLE IF EXISTS public.environment_permissions;

-- environment_user_secrets
ALTER TABLE public.environment_user_secrets DROP CONSTRAINT fk_environment_user_secrets_secret;
ALTER TABLE public.environment_user_secrets DROP CONSTRAINT fk_environment_user_secrets_user;
ALTER TABLE public.environment_user_secrets DROP CONSTRAINT fk_environment_user_secrets_environment;
DROP TABLE IF EXISTS public.environment_user_secrets;

-- project_environment_secrets
/* ALTER TABLE public.project_environment_secrets DROP CONSTRAINT fk_project_environment_secrets_secret; */
/* ALTER TABLE public.project_environment_secrets DROP CONSTRAINT fk_project_environment_secrets_environment; */
/* DROP TABLE IF EXISTS public.project_environment_secrets; */

-- secrets
DROP TABLE IF EXISTS public.secrets;
