INSERT INTO
  public.environment_types (id, name, created_at, updated_at)
VALUES
  (1, 'dev', current_timestamp, current_timestamp),
  (2, 'staging', current_timestamp, current_timestamp),
  (3, 'ci', current_timestamp, current_timestamp),
  (4, 'prod', current_timestamp, current_timestamp);

INSERT INTO
  public.roles (id, name, description, created_at, updated_at)
VALUES
  (
    1,
    'developer',
    'Can read and write the development environment',
    current_timestamp,
    current_timestamp
  ),
  (
    2,
    'developer (invite)',
    'Same as developer, and can add new members',
    current_timestamp,
    current_timestamp
  ),
  (
    3,
    'devops',
    'Can read, write and add new members on staging, ci and prod',
    current_timestamp,
    current_timestamp
  ),
  (
    4,
    'admin',
    'Can read, write and add new members on all environments',
    current_timestamp,
    current_timestamp
  );

INSERT INTO
  public.roles_environment_types (id, role_id, environment_type_id, read, write, invite, created_at, updated_at)
VALUES
  -- developer
  (1,  1, 1, true, true, false, current_timestamp, current_timestamp),
  (2,  1, 2, false, false, false, current_timestamp, current_timestamp),
  (3,  1, 3, false, false, false, current_timestamp, current_timestamp),
  (4,  1, 4, true, true, false, current_timestamp, current_timestamp),
  -- developer (invite)
  (5,  2, 1, true, true, true, current_timestamp, current_timestamp),
  (6,  2, 2, false, false, false, current_timestamp, current_timestamp),
  (7,  2, 3, false, false, false, current_timestamp, current_timestamp),
  (8,  2, 4, false, false, false, current_timestamp, current_timestamp),
  -- devops
  (9,  3, 1, false, false, false, current_timestamp, current_timestamp),
  (10, 3, 2, true, true, true, current_timestamp, current_timestamp),
  (11, 3, 3, true, true, true, current_timestamp, current_timestamp),
  (12, 3, 4, true, true, true, current_timestamp, current_timestamp),
  -- admin
  (13, 4, 1, true, true, true, current_timestamp, current_timestamp),
  (14, 4, 2, true, true, true, current_timestamp, current_timestamp),
  (15, 4, 3, true, true, true, current_timestamp, current_timestamp),
  (16, 4, 4, true, true, true, current_timestamp, current_timestamp);

