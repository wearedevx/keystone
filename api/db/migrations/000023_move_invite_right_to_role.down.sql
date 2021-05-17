ALTER TABLE
  public.roles_environment_types
ADD
  COLUMN invite boolean NOT NULL DEFAULT FALSE;
--
ALTER TABLE
  public.roles DROP COLUMN can_add_member;
-- =================
-- SEEDS
-- =================
-- developer cannot add
UPDATE
  public.roles_environment_types
SET
  invite = TRUE
WHERE
  id IN (5, 9, 13, 14, 16);
