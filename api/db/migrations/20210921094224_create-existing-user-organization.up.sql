INSERT INTO
  public.organizations (owner_id, name, created_at, updated_at)
SELECT
  id as owner_id,
  user_id as name,
  created_at,
  updated_at
FROM
  public.users
