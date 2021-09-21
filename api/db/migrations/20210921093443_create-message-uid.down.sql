ALTER TABLE
    public.messages DROP COLUMN uuid;

ALTER TABLE
    public.messages
ALTER COLUMN
    payload
SET
    NOT NULL;
