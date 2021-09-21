ALTER TABLE
  public.roles
ADD
  COLUMN can_add_member boolean NOT NULL DEFAULT FALSE;
ALTER TABLE
  public.roles_environment_types DROP COLUMN invite;
-- =================
-- SEEDS
-- =================
-- developer cannot add
UPDATE
  public.roles
SET
  can_add_member = FALSE
WHERE
  name = 'developer';
-- developer (invite) can add
UPDATE
  public.roles
SET
  can_add_member = TRUE
WHERE
  name = 'developer (invite)';
-- devops can add
UPDATE
  public.roles
SET
  can_add_member = TRUE
WHERE
  name = 'devops';
-- admin can add
UPDATE
  public.roles
SET
  can_add_member = TRUE
WHERE
  name = 'admin';
