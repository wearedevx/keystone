ALTER TABLE
  public.devices
DROP
  CONSTRAINT fk_users_public_keys;

ALTER TABLE
  public.devices
ALTER COLUMN user_id drop not NULL;
