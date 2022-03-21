DROP INDEX IF EXISTS idx_public_keys_user_id;

ALTER TABLE public.devices
DROP COLUMN IF EXISTS user_id;
