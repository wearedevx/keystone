CREATE TABLE IF NOT EXISTS public.projects (
	id bigserial NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	deleted_at timestamptz NULL,
	uuid text NOT NULL,
	"name" text NOT NULL,
	CONSTRAINT projects_pkey PRIMARY KEY (id),
	CONSTRAINT projects_uuid_key UNIQUE (uuid)
);
CREATE INDEX idx_projects_deleted_at ON public.projects USING btree (deleted_at);
