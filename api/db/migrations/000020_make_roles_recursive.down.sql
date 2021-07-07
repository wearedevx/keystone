ALTER TABLE
  public.roles DROP CONSTRAINT fk_roles_parent;
DROP INDEX IF EXISTS idx_roles_parent_id;
ALTER TABLE
  public.roles DROP COLUMN parent_id integer;
