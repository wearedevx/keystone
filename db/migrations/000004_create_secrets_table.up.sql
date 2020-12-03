CREATE TABLE IF NOT EXISTS public.secrets (
	id bigserial NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	deleted_at timestamptz NULL,
	"name" text NOT NULL,
	"type" text NOT NULL,
	CONSTRAINT secrets_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_secrets_deleted_at ON public.secrets USING btree (deleted_at);
