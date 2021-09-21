ALTER TABLE
  devices
ADD
  COLUMN IF NOT EXISTS last_used_at timestamptz NULL;
