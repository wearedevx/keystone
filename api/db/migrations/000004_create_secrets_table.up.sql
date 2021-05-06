CREATE TABLE IF NOT EXISTS public.secrets (
	id bigserial NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	"name" text NOT NULL,
	"type" text NOT NULL,
	CONSTRAINT secrets_pkey PRIMARY KEY (id)
);

