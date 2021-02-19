CREATE TABLE IF NOT EXISTS public.messages (
	id bigserial NOT NULL,
	sender_id integer NOT NULL,
	recipient_id integer NOT NULL,
	payload bytea NOT NULL,
	created_at timestamptz NOT NULL,
	updated_at timestamptz NOT NULL,
	CONSTRAINT messages_pkey PRIMARY KEY (id)
);

DROP INDEX IF EXISTS idx_messages_sender_id;
CREATE INDEX idx_messages_sender_id ON public.messages USING btree (sender_id);
DROP INDEX IF EXISTS idx_messages_recipient_id;
CREATE INDEX idx_messages_recipient_id ON public.messages USING btree (recipient_id);

-- messages foreign keys

ALTER TABLE public.messages ADD CONSTRAINT fk_messages_sender FOREIGN KEY (sender_id) REFERENCES public.project_members ON DELETE CASCADE ON UPDATE NO ACTION;
ALTER TABLE public.messages ADD CONSTRAINT fk_messages_recipient FOREIGN KEY (recipient_id) REFERENCES public.project_members ON DELETE CASCADE ON UPDATE NO ACTION;
