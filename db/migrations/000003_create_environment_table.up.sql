CREATE TABLE IF NOT EXISTS public.environments (
	id bigserial NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	deleted_at timestamptz NULL,
	"name" text NOT NULL,
	project_id int8 NOT NULL,
	CONSTRAINT environments_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_environments_deleted_at ON public.environments USING btree (deleted_at);
CREATE INDEX idx_environments_project_id ON public.environments USING btree (project_id);


-- public.environments foreign keys

ALTER TABLE public.environments ADD CONSTRAINT fk_projects_environments FOREIGN KEY (project_id) REFERENCES projects(id);
