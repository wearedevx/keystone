DELETE FROM public_keys;

ALTER TABLE
  public_keys
ADD
  COLUMN uid text NOT NULL;


ALTER TABLE public_keys RENAME TO devices;

ALTER TABLE devices RENAME COLUMN device TO name;

ALTER TABLE devices RENAME COLUMN key TO public_key;
