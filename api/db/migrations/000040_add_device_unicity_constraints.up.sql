DROP INDEX IF EXISTS idx_devices_public_key;
CREATE UNIQUE INDEX idx_devices_public_key ON public.devices USING btree (public_key);

DROP INDEX IF EXISTS idx_user_devices_user_id_device_id;
CREATE UNIQUE INDEX idx_user_devices_user_id_device_id ON public.user_devices USING btree (user_id, device_id);
