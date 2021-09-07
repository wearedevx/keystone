CREATE TABLE IF NOT EXISTS public.user_devices (
  id bigserial NOT NULL,
  user_id bigserial NOT NULL,
  device_id bigserial NOT NULL,
  created_at timestamptz NULL,
  updated_at timestamptz NULL,
  CONSTRAINT user_devices_pkey PRIMARY KEY (id)
);


DROP INDEX IF EXISTS idx_user_devices_user_id;
CREATE INDEX idx_user_devices_user_id ON public.user_devices USING btree (user_id);
ALTER TABLE
  public.user_devices
ADD
  CONSTRAINT fk_user_devices_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON
UPDATE
  NO ACTION;



DROP INDEX IF EXISTS idx_user_devices_device_id;
CREATE INDEX idx_user_devices_device_id ON public.user_devices USING btree (device_id);
ALTER TABLE
  public.user_devices
ADD
  CONSTRAINT fk_user_devices_devices FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE ON
UPDATE
  NO ACTION;



/* ALTER TABLE public.users ADD COLUMN user_devices_id bigserial NOT NULL; */

/* ALTER TABLE public.devices ADD COLUMN user_devices_id bigserial NOT NULL; */

/* DROP INDEX IF EXISTS idx_devices_user_devices_id; */
/* CREATE UNIQUE INDEX idx_devices_user_devices_id ON public.devices USING btree (user_devices_id); */

/* DROP INDEX IF EXISTS idx_users_user_devices_id; */
/* CREATE UNIQUE INDEX idx_users_user_devices_id ON public.users USING btree (user_devices_id); */

/* ALTER TABLE */
/*   public.devices */
/* ADD */
/*   CONSTRAINT fk_devices_user_devices FOREIGN KEY (user_devices_id) REFERENCES user_devices(id) ON DELETE CASCADE ON */
/* UPDATE */
/*   NO ACTION; */


/* ALTER TABLE */
/*   public.users */
/* ADD */
/*   CONSTRAINT fk_users_user_devices FOREIGN KEY (user_devices_id) REFERENCES user_devices(id) ON DELETE CASCADE ON */
/* UPDATE */
/*   NO ACTION; */


