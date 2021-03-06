CREATE TABLE IF NOT EXISTS public.users (
	id bigserial NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	account_type text NULL DEFAULT 'custom'::text,
	user_id text NULL,
	ext_id text NULL,
	username text NULL,
	fullname text NOT NULL,
	email text NOT NULL,
	keys_cipher text NULL,
	keys_sign text NULL,
	CONSTRAINT users_pkey PRIMARY KEY (id)
);

DROP INDEX IF EXISTS idx_users_user_id;
CREATE UNIQUE INDEX idx_users_user_id ON public.users USING btree (user_id);

DROP INDEX IF EXISTS idx_users_username;
CREATE UNIQUE INDEX idx_users_username ON public.users USING btree (username);
