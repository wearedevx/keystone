CREATE UNIQUE INDEX idx_environments_environment_id ON public.environments USING btree (environment_id);

ALTER TABLE public.messages ADD COLUMN environment_id text NOT NULL;

ALTER TABLE public.messages ADD CONSTRAINT fk_messages_environments FOREIGN KEY (environment_id) REFERENCES environments(environment_id);
