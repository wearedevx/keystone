INSERT INTO
  public.roles (id, name, description, created_at, updated_at)
VALUES
  (
    5,
    'ci',
    'Can read on all environments',
    current_timestamp,
    current_timestamp
  );

INSERT INTO
  public.roles_environment_types (id, role_id, environment_type_id, read, write, created_at, updated_at)
VALUES
  (17, 5, 1, true, false, current_timestamp, current_timestamp),
  (18, 5, 2, true, false, current_timestamp, current_timestamp),
  (19, 5, 3, true, false, current_timestamp, current_timestamp);
