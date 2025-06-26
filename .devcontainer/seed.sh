DB=postgres://postgres:postgres@localhost:5432/postgres

# Wait for backend to be ready
until [ "$(curl -s -m 1 http://localhost:5002/tribes 2>/dev/null)" = "[]" ]
do
  echo "Waiting for backend to become ready..."
  sleep 5
done

echo "Inserting dummy data...."

count=$(psql $DB -tA -c "SELECT count(*) FROM people")
if [ "$count" -gt 0 ]; then
  echo "Dummy data exists!"
else
  psql $DB -f docker/dummy-data/people.sql
  psql $DB -f docker/dummy-data/workspaces.sql
  psql $DB -f docker/dummy-data/paid-bounties.sql
fi



