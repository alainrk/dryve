#!/bin/bash

set -eu

echo "Creating necessary directories..."
mkdir -p .data .database


echo "Setting up the database..."
docker-compose up -d --build db

sleep 5

echo "Running migrations..."
make automigrate

echo "Setting up remaining components..."
docker-compose up -d

exit 0
