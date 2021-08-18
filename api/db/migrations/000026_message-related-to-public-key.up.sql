DELETE FROM public.messages;

ALTER TABLE public.messages
ADD COLUMN public_key_id integer NOT NULL;

DROP INDEX IF EXISTS idx_messages_public_key_id;
CREATE INDEX idx_messages_public_key_id ON public.messages USING btree (public_key_id);

ALTER TABLE public.messages
ADD CONSTRAINT fk_messages_public_key FOREIGN KEY (public_key_id) REFERENCES public.public_keys(id);
