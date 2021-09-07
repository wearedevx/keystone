CREATE TABLE IF NOT EXISTS public.user_device (
  id bigserial NOT NULL,
  user_id bigserial NOT NULL,
  device_id bigserial NOT NULL,
  created_at timestamptz NULL,
  updated_at timestamptz NULL,
  CONSTRAINT user_device_pkey PRIMARY KEY (id)
);


DROP INDEX IF EXISTS idx_user_device_user_id;
CREATE INDEX idx_user_device_user_id ON public.user_device USING btree (user_id);
ALTER TABLE
  public.user_device
ADD
  CONSTRAINT fk_user_device_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON
UPDATE
  NO ACTION;



DROP INDEX IF EXISTS idx_user_device_device_id;
CREATE INDEX idx_user_device_device_id ON public.user_device USING btree (device_id);
ALTER TABLE
  public.user_device
ADD
  CONSTRAINT fk_user_device_devices FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE ON
UPDATE
  NO ACTION;



/* ALTER TABLE public.users ADD COLUMN user_device_id bigserial NOT NULL; */

/* ALTER TABLE public.devices ADD COLUMN user_device_id bigserial NOT NULL; */

/* DROP INDEX IF EXISTS idx_devices_user_device_id; */
/* CREATE UNIQUE INDEX idx_devices_user_device_id ON public.devices USING btree (user_device_id); */

/* DROP INDEX IF EXISTS idx_users_user_device_id; */
/* CREATE UNIQUE INDEX idx_users_user_device_id ON public.users USING btree (user_device_id); */

/* ALTER TABLE */
/*   public.devices */
/* ADD */
/*   CONSTRAINT fk_devices_user_device FOREIGN KEY (user_device_id) REFERENCES user_device(id) ON DELETE CASCADE ON */
/* UPDATE */
/*   NO ACTION; */


/* ALTER TABLE */
/*   public.users */
/* ADD */
/*   CONSTRAINT fk_users_user_device FOREIGN KEY (user_device_id) REFERENCES user_device(id) ON DELETE CASCADE ON */
/* UPDATE */
/*   NO ACTION; */


