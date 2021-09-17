CREATE TABLE IF NOT EXISTS public.public_keys (
  id bigserial NOT NULL,
  user_id bigserial NOT NULL,
  key bytea NOT NULL,
  created_at timestamptz NULL,
  updated_at timestamptz NULL,
  CONSTRAINT public_keys_pkey PRIMARY KEY (id)
);

DROP INDEX IF EXISTS idx_public_keys_user_id;
CREATE INDEX idx_public_keys_user_id ON public.public_keys USING btree (user_id);
ALTER TABLE
  public.public_keys
ADD
  CONSTRAINT fk_users_public_keys FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON
UPDATE
  NO ACTION;
;

INSERT INTO
  public.public_keys (user_id, key, created_at, updated_at)
SELECT
  id as user_id,
  public_key as key,
  created_at,
  updated_at
FROM
  public.users
