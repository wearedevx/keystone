CREATE TABLE IF NOT EXISTS public.messages (
	id bigserial NOT NULL,
	sender_id integer NOT NULL,
	recipient_id integer NOT NULL,
	payload bytea NOT NULL,
	created_at timestamptz NOT NULL,
	updated_at timestamptz NOT NULL,
	CONSTRAINT messages_pkey PRIMARY KEY (id)
);

-- messages foreign keys

ALTER TABLE public.messages ADD CONSTRAINT fk_messages_sender FOREIGN KEY (sender_id) REFERENCES public.project_members ON DELETE CASCADE ON UPDATE NO ACTION;
ALTER TABLE public.messages ADD CONSTRAINT fk_messages_recipient FOREIGN KEY (recipient_id) REFERENCES public.project_members ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE public.messages DROP CONSTRAINT fk_messages_sender;
ALTER TABLE public.messages DROP CONSTRAINT fk_messages_recipient;

DROP TABLE IF EXISTS public.messages;
