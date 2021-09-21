ALTER TABLE
    public.messages
ADD
    COLUMN uuid TEXT;

ALTER TABLE
    public.messages
ALTER COLUMN
    payload DROP NOT NULL;
