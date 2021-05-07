ALTER TABLE public.messages DROP CONSTRAINT fk_messages_sender;
ALTER TABLE public.messages DROP CONSTRAINT fk_messages_recipient;

ALTER TABLE public.messages ADD CONSTRAINT fk_messages_sender FOREIGN KEY (sender_id) REFERENCES public.users ON DELETE CASCADE ON UPDATE NO ACTION;
ALTER TABLE public.messages ADD CONSTRAINT fk_messages_recipient FOREIGN KEY (recipient_id) REFERENCES public.users ON DELETE CASCADE ON UPDATE NO ACTION;
