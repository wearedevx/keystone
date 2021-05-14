ALTER TABLE
  public.roles
ADD
  COLUMN parent_id integer;
DROP INDEX IF EXISTS idx_roles_parent_id;
CREATE INDEX idx_roles_parent_id ON public.roles USING btree (parent_id);
ALTER TABLE
  public.roles DROP CONSTRAINT IF EXISTS fk_roles_parent;
ALTER TABLE
  public.roles
ADD
  CONSTRAINT fk_roles_parent FOREIGN KEY (parent_id) REFERENCES public.roles(id);
