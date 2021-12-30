ALTER TABLE public.devices
  ADD COLUMN IF NOT EXISTS newly_created boolean;
