alter table public.projects add column organization_id integer;


DROP INDEX IF EXISTS idx_projects_organization_id;
CREATE INDEX idx_projects_organization_id ON public.projects USING btree (organization_id);

UPDATE public.projects p
SET organization_id =  ( SELECT id FROM public.organizations o WHERE o.owner_id  = p.user_id);

alter table public.projects
alter column organization_id type integer,
alter column organization_id set not null,
ADD CONSTRAINT fk_project_organization FOREIGN KEY (organization_id) REFERENCES public.organizations(id);

