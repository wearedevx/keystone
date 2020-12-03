CREATE TABLE public.environment_permissions (
	user_id int8 NOT NULL,
	environment_id int8 NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	"role" text NULL,
	CONSTRAINT environment_permissions_pkey PRIMARY KEY (user_id, environment_id)
);


-- public.environment_permissions foreign keys

ALTER TABLE public.environment_permissions ADD CONSTRAINT fk_environment_permissions_environment FOREIGN KEY (environment_id) REFERENCES environments(id);
ALTER TABLE public.environment_permissions ADD CONSTRAINT fk_environment_permissions_user FOREIGN KEY (user_id) REFERENCES users(id);
