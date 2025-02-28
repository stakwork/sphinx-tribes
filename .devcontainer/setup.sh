#!/bin/bash

cd /workspaces

git clone https://github.com/stakwork/sphinx-tribes-frontend 

cd sphinx-tribes-frontend

yarn

# Start the frontend service in the background
echo "Starting frontend service..."
yarn start:codespace &
FRONTEND_PID=$!

# Setup a trap to kill the frontend service when script exits
trap 'echo "Shutting down frontend service..."; kill $FRONTEND_PID' EXIT

cd /workspaces/sphinx-tribes

go build

# Start the backend service
echo "Starting backend service..."
./sphinx-tribes