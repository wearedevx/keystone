-- devops
UPDATE
  public.roles
SET
  parent_id = (
    SELECT
      id
    FROM
      public.roles
    WHERE
      name = 'admin'
  )
WHERE
  name = 'devops';
-- developer (invite)
UPDATE
  public.roles
SET
  parent_id = (
    SELECT
      id
    FROM
      public.roles
    WHERE
      name = 'devops'
  )
WHERE
  name = 'developer (invite)';
-- developer
UPDATE
  public.roles
SET
  parent_id = (
    SELECT
      id
    FROM
      public.roles
    WHERE
      name = 'developer (invite)'
  )
WHERE
  name = 'developer';
