-- devops
UPDATE
  public.roles
SET
  parent_id = NULL
WHERE
  name = 'devops';
-- developer (invite)
UPDATE
  public.roles
SET
  parent_id = NULL
WHERE
  name = 'developer (invite)';
-- developer
UPDATE
  public.roles
SET
  parent_id = NULL
WHERE
  name = 'developer';
