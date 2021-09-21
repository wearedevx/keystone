truncate table messages;

ALTER TABLE public.messages
DROP CONSTRAINT fk_messages_recipient_device_id;

alter table messages rename column  recipient_device_id to public_key_id;

ALTER TABLE public.messages
ADD CONSTRAINT fk_messages_public_key_id FOREIGN KEY (public_keys) REFERENCES public.public_keys(id);

alter table messages drop column sender_device_id;
