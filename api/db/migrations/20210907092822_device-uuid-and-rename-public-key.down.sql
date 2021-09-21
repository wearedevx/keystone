ALTER TABLE
  public_keys drop column uid

alter table devices rename to public_keys;

alter table public_keys rename column name to device;

alter table public_keys rename column public_key to key;
