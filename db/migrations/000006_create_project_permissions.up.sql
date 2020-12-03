CREATE TABLE public.project_permissions (
	user_id int8 NOT NULL,
	project_id int8 NOT NULL,
	"role" text NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	deleted_at timestamptz NULL,
	CONSTRAINT project_permissions_pkey PRIMARY KEY (user_id, project_id)
);


-- public.project_permissions foreign keys

ALTER TABLE public.project_permissions ADD CONSTRAINT fk_project_permissions_project FOREIGN KEY (project_id) REFERENCES projects(id);
ALTER TABLE public.project_permissions ADD CONSTRAINT fk_project_permissions_user FOREIGN KEY (user_id) REFERENCES users(id);
