ALTER TABLE projects
ADD COLUMN ttl INTEGER NOT NULL DEFAULT 7,
ADD COLUMN days_before_ttl_expiry INTEGER NOT NULL DEFAULT 2;
