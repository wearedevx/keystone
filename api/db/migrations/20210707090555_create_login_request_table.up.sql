CREATE TABLE IF NOT EXISTS public.login_requests (
	id bigserial NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	temporary_code text NOT NULL,
	auth_code text NULL,
	answered bool NULL DEFAULT false,
	CONSTRAINT login_requests_pkey PRIMARY KEY (id)
);

