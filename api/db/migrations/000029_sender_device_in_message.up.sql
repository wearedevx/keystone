truncate table messages;

ALTER TABLE public.messages
DROP CONSTRAINT fk_messages_public_key;

alter table messages rename column public_key_id to recipient_device_id;

ALTER TABLE public.messages
ADD CONSTRAINT fk_messages_recipient_device_id FOREIGN KEY (recipient_device_id) REFERENCES public.devices(id);

alter table messages ADD column sender_device_id int not null;

ALTER TABLE public.messages
ADD CONSTRAINT fk_messages_sender_device_id FOREIGN KEY (sender_device_id) REFERENCES public.devices(id);
