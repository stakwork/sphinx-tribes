
-- TRIBES

CREATE TABLE tribes (
  uuid TEXT NOT NULL PRIMARY KEY,
  owner_pub_key TEXT NOT NULL,
  owner_alias TEXT,
  group_key TEXT,
  name TEXT,
  description TEXT,
  tags TEXT[] not null default '{}',
  img TEXT,
  price_to_join BIGINT,
  price_per_message BIGINT,
  escrow_amount BIGINT,
  escrow_millis BIGINT,
  created timestamptz,
  updated timestamptz,
  member_count BIGINT,
  unlisted boolean,
  private boolean,
  deleted boolean,
  app_url TEXT,
  feed_url TEXT,
  second_brain_url TEXT,
  feed_type INT,
  last_active BIGINT,
  bots TEXT,
  owner_route_hint TEXT,
  unique_name TEXT,
  pin TEXT,
  preview TEXT,
  profile_filters TEXT,
  second_brain_url TEXT
);

-- for searching 

ALTER TABLE tribes ADD COLUMN tsv tsvector;

UPDATE tribes SET tsv =
  setweight(to_tsvector(name), 'A') ||
	setweight(to_tsvector(description), 'B') ||
	setweight(array_to_tsvector(tags), 'C');

CREATE INDEX tribes_tsv ON tribes USING GIN(tsv);

-- select

SELECT name, description, tags
FROM tribes, to_tsquery('foo') q
WHERE tsv @@ q;

-- rank

SELECT name, uuid, description, ts_rank(tsv, q) as rank
FROM tribes, to_tsquery('anothe') q
WHERE tsv @@ q
ORDER BY rank DESC
LIMIT 12;

-- plainto_tsquery is another way

-- BOTS
CREATE TABLE bots (
  uuid TEXT NOT NULL PRIMARY KEY,
  owner_pub_key TEXT NOT NULL,
  owner_alias TEXT,
  name TEXT,
  unique_name TEXT,
  description TEXT,
  tags TEXT[] not null default '{}',
  img TEXT,
  price_per_use BIGINT,
  created timestamptz,
  updated timestamptz,
  member_count BIGINT,
  unlisted boolean,
  deleted boolean
);

-- for searching 

ALTER TABLE bots ADD COLUMN tsv tsvector;

UPDATE bots SET tsv =
  setweight(to_tsvector(name), 'A') ||
	setweight(to_tsvector(description), 'B') ||
	setweight(array_to_tsvector(tags), 'C');

CREATE INDEX bots_tsv ON bots USING GIN(tsv);

SELECT uuid, unique_name, ts_rank(tsv, q) as rank
  FROM bots, to_tsquery('btc') q
  WHERE tsv @@ q
  ORDER BY rank DESC LIMIT 2 OFFSET 0;

-- PEOPLE

CREATE TABLE people (
  id SERIAL PRIMARY KEY,
  uuid TEXT,
  owner_pub_key TEXT NOT NULL,
  owner_alias TEXT,
  owner_route_hint TEXT,
  owner_contact_key TEXT,
  description TEXT,
  tags TEXT[] not null default '{}',
  img TEXT,
  created timestamptz,
  updated timestamptz,
  unlisted boolean,
  deleted boolean,
  unique_name TEXT,
  price_to_meet BIGINT,
  extras JSONB,
  twitter_confirmed BOOLEAN,
  github_issues JSONB
);

ALTER TABLE people ADD COLUMN tsv tsvector;

UPDATE people SET tsv =
  setweight(to_tsvector(owner_alias), 'A') ||
	setweight(to_tsvector(description), 'B') ||
	setweight(array_to_tsvector(tags), 'C');

CREATE INDEX people_tsv ON people USING GIN(tsv);

INSERT into people (owner_alias, owner_pub_key, description, tags, img, unique_name)
VALUES
('Evan', '02290714deafd0cb33d2be3b634fc977a98a9c9fa1dd6c53cf17d99b350c08c67b', 'Im cool', '{"tag1"}', 'https://evan.cool/img/trumpetplay.jpg', 'evan');

INSERT into people (owner_alias, owner_pub_key, description, tags, img, unique_name)
VALUES
('Jesse', '038c3c1f4d304c7b997fecfdaf8fdfc2215405942c025349b45de9dfe6fdb8a43e', 'Im cool', '{"tag1"}', 'https://cliparting.com/wp-content/uploads/2018/03/cool-pictures-2018-2.jpg', 'jesse');

ALTER TABLE IF EXISTS tribes ADD COLUMN IF NOT EXISTS preview VARCHAR NULL;

CREATE TABLE connectioncodes (
  id SERIAL PRIMARY KEY,
  connection_string TEXT,
  is_used boolean,
  date_created timestamptz,
  pubkey TEXT,
  route_hint TEXT,
  sats_amount bigint
)