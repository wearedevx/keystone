CREATE TABLE IF NOT EXISTS public.activity_logs (
	id bigserial NOT NULL,
	user_id INTEGER NULL,
	project_id INTEGER NULL,
	environment_id INTEGER NULL,
	action VARCHAR(255) DEFAULT '',
	success BOOLEAN NOT NULL DEFAULT FALSE,
	error VARCHAR(255) DEFAULT '',
	created_at TIMESTAMPTZ NOT NULL,
	updated_at TIMESTAMPTZ,
	CONSTRAINT acitivity_logs_pkey PRIMARY KEY (id)
);

DROP INDEX IF EXISTS idx_activity_logs_user_id;
CREATE INDEX idx_activity_logs_user_id ON public.activity_logs USING btree (user_id);

DROP INDEX IF EXISTS idx_activity_logs_project_id;
CREATE INDEX idx_activity_logs_project_id ON public.activity_logs USING btree (project_id);

DROP INDEX IF EXISTS idx_activity_logs_environment_id;
CREATE INDEX idx_activity_logs_environment_id ON public.activity_logs USING btree (environment_id);

-- public.activity_logs foreign keys

ALTER TABLE public.activity_logs ADD CONSTRAINT fk_activity_logs_project FOREIGN KEY (project_id) REFERENCES projects(id);
ALTER TABLE public.activity_logs ADD CONSTRAINT fk_activity_logs_user FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE public.activity_logs ADD CONSTRAINT fk_activity_logs_environment FOREIGN KEY (environment_id) REFERENCES environments(id);
