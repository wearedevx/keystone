
ALTER TABLE public.environments ADD COLUMN project_id integer NOT NULL;
DROP INDEX IF EXISTS idx_environments_project_id;
CREATE INDEX idx_environments_project_id ON public.environments USING btree (project_id);
ALTER TABLE public.environments ADD CONSTRAINT fk_projects_environments FOREIGN KEY (project_id) REFERENCES projects(id);

ALTER TABLE public.users ADD COLUMN keys_cipher text NULL;
ALTER TABLE public.users ADD COLUMN keys_sign text NULL;
