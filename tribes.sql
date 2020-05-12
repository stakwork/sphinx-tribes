
CREATE TABLE tribes (
  uuid TEXT NOT NULL PRIMARY KEY,
  owner_pub_key TEXT NOT NULL,
  group_key TEXT,
  name TEXT,
  description TEXT,
  tags TEXT[] not null default '{}',
  img TEXT,
  price_to_join BIGINT,
  price_per_message BIGINT,
  created timestamptz,
  updated timestamptz
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


