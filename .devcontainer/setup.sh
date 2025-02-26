#!/bin/bash

echo "Starting PostgreSQL..."
sudo service postgresql start

echo "Creating database and user..."
sudo -u postgres psql <<EOF
CREATE DATABASE sphinx_tribes;
CREATE USER sphinx_user WITH ENCRYPTED PASSWORD 'yourpassword';
GRANT ALL PRIVILEGES ON DATABASE sphinx_tribes TO sphinx_user;
EOF

echo "Setup complete!"
