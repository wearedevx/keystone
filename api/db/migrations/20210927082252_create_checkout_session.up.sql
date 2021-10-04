CREATE TABLE IF NOT EXISTS public.checkout_sessions (
	id bigserial NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	session_id varchar(255) NULL,
	status varchar(12) DEFAULT 'pending'
	);

ALTER TABLE public.checkout_sessions
	DROP CONSTRAINT IF EXISTS check_checkout_status;

ALTER TABLE public.checkout_sessions
	ADD CONSTRAINT check_checkout_status CHECK (status IN ('pending', 'success', 'canceled'));
