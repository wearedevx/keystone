ALTER TABLE public.organizations
  ADD COLUMN IF NOT EXISTS customer_id varchar(255);
