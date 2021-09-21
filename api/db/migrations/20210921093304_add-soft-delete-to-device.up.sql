ALTER TABLE
  devices
ADD
  COLUMN IF NOT EXISTS deleted_at timestamptz NULL;
