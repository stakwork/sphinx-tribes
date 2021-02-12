
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
  last_active timestamptz,
  bots TEXT,
  owner_route_hint TEXT,
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

SELECT name, id, description, ts_rank(tsv, q) as rank
FROM tribes, to_tsquery('anothe') q
WHERE tsv @@ q
ORDER BY rank DESC
LIMIT 12;

-- plainto_tsquery is another way






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