ALTER TABLE
  public.devices
ADD
  CONSTRAINT fk_users_public_keys FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON
UPDATE
  NO ACTION;
;

ALTER TABLE
  public.devices
ALTER COLUMN user_id set not NULL;
