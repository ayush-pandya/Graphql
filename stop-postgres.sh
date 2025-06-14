#!/bin/bash

# Script to stop PostgreSQL container

CONTAINER_NAME="ticketdb-postgres"

echo "🛑 Stopping PostgreSQL container..."

if [ "$(docker ps -q -f name=$CONTAINER_NAME)" ]; then
    docker stop $CONTAINER_NAME
    echo "✅ PostgreSQL container stopped"
else
    echo "ℹ️  Container is not running"
fi

# Uncomment the line below if you want to remove the container completely
# docker rm $CONTAINER_NAME 