CREATE TABLE IF NOT EXISTS public.organizations (
	id bigserial NOT NULL,
	owner_id integer NOT NULL,
	name text NOT NULL,
	paid boolean NOT NULL DEFAULT false,
	created_at timestamptz NOT NULL,
	updated_at timestamptz NOT NULL,
	CONSTRAINT organizations_pkey PRIMARY KEY (id)
);

DROP INDEX IF EXISTS idx_organizations_owner_id;
CREATE INDEX idx_organizations_owner_id ON public.organizations USING btree (owner_id);

-- organizations foreign keys

ALTER TABLE public.organizations ADD CONSTRAINT fk_organization_owner FOREIGN KEY (owner_id) REFERENCES public.users ON DELETE CASCADE ON UPDATE NO ACTION;
