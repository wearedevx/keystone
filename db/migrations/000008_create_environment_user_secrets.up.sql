CREATE TABLE public.environment_user_secrets (
	user_id int8 NOT NULL,
	secret_id int8 NOT NULL,
	environment_id int8 NULL,
	value bytea NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	deleted_at timestamptz NULL,
	CONSTRAINT environment_user_secrets_pkey PRIMARY KEY (user_id, secret_id, environment_id)
);


-- public.environment_user_secrets foreign keys

ALTER TABLE public.environment_user_secrets ADD CONSTRAINT fk_environment_user_secrets_secret FOREIGN KEY (secret_id) REFERENCES secrets(id);
ALTER TABLE public.environment_user_secrets ADD CONSTRAINT fk_environment_user_secrets_user FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE public.environment_user_secrets ADD CONSTRAINT fk_environment_user_secrets_environment FOREIGN KEY (environment_id) REFERENCES environments(id);
