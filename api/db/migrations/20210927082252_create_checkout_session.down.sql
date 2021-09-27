ALTER TABLE public.checkout_sessions
	DROP CONSTRAINT IF EXISTS check_checkout_status;

DROP TABLE IF EXISTS public.checkout_sessions;
