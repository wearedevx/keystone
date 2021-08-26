TRUNCATE TABLE messages;

ALTER TABLE public.messages
DROP CONSTRAINT fk_messages_public_key;

ALTER TABLE messages RENAME COLUMN public_key_id TO recipient_device_id;

ALTER TABLE public.messages
ADD CONSTRAINT fk_messages_recipient_device_id FOREIGN KEY (recipient_device_id) REFERENCES public.devices(id);

ALTER TABLE messages ADD COLUMN sender_device_id int NOT NULL;

ALTER TABLE public.messages
ADD CONSTRAINT fk_messages_sender_device_id FOREIGN KEY (sender_device_id) REFERENCES public.devices(id);
