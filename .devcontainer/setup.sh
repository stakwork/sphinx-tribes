#!/bin/bash

cd /workspaces

git clone https://github.com/stakwork/sphinx-tribes-frontend

cd sphinx-tribes

DB=postgres://postgres:postgres@localhost:5432/postgres

until psql $DB -c '\q'
do
  echo "Waiting for PostgreSQL to become available..."
  sleep 5
done

echo "Inserting dummy data...."

psql $DB -f docker/dummy-data/people.sql
psql $DB -f docker/dummy-data/paid-bounties.sql

gh codespace ports visibility 5002:public -c $CODESPACE_NAME
gh codespace ports visibility 13008:public -c $CODESPACE_NAME
gh codespace ports visibility 15552:public -c $CODESPACE_NAME