ALTER TABLE
  public_keys
add
  column uid text not null;


alter table public_keys rename to devices;

alter table devices rename column device to name;

alter table devices rename column key to public_key;
