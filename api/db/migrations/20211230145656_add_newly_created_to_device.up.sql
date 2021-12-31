ALTER TABLE public.user_devices
  ADD COLUMN IF NOT EXISTS newly_created boolean;
