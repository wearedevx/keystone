ALTER TABLE public.activity_logs
DROP CONSTRAINT fk_activity_logs_project;

ALTER TABLE public.activity_logs
DROP CONSTRAINT fk_activity_logs_user;

ALTER TABLE public.activity_logs
DROP CONSTRAINT fk_activity_logs_environment;

DROP INDEX IF EXISTS idx_activity_logs_user_id;
DROP INDEX IF EXISTS idx_activity_logs_project_id;
DROP INDEX IF EXISTS idx_activity_logs_environment_id;

DROP TABLE IF EXISTS public.activity_logs;
