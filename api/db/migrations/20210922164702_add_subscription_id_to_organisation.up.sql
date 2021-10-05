ALTER TABLE public.organizations
  ADD COLUMN IF NOT EXISTS subscription_id VARCHAR(255);
