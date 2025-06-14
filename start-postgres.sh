#!/bin/bash

# Script to start PostgreSQL using Docker
# This will create and start a PostgreSQL container for your Go application

CONTAINER_NAME="ticketdb-postgres"
DB_NAME="ticketdb"
DB_USER="ayushpandya"
DB_PASSWORD="postgres"
DB_PORT="5432"

echo "🐘 Starting PostgreSQL container..."

# Check if container already exists
if [ "$(docker ps -aq -f name=$CONTAINER_NAME)" ]; then
    echo "📦 Container $CONTAINER_NAME already exists"
    
    # Check if it's running
    if [ "$(docker ps -q -f name=$CONTAINER_NAME)" ]; then
        echo "✅ Container is already running"
    else
        echo "🔄 Starting existing container..."
        docker start $CONTAINER_NAME
    fi
else
    echo "🆕 Creating new PostgreSQL container..."
    docker run -d \
        --name $CONTAINER_NAME \
        -e POSTGRES_DB=$DB_NAME \
        -e POSTGRES_USER=$DB_USER \
        -e POSTGRES_PASSWORD=$DB_PASSWORD \
        -p $DB_PORT:5432 \
        -v postgres_data:/var/lib/postgresql/data \
        postgres:15
fi

# Wait for PostgreSQL to be ready
echo "⏳ Waiting for PostgreSQL to be ready..."
sleep 5

# Test connection
docker exec $CONTAINER_NAME pg_isready -U $DB_USER -d $DB_NAME

if [ $? -eq 0 ]; then
    echo "✅ PostgreSQL is running and ready!"
    echo "📊 Database: $DB_NAME"
    echo "👤 User: $DB_USER"
    echo "🔌 Port: $DB_PORT"
    echo ""
    echo "🔗 Connection string:"
    echo "postgres://$DB_USER:$DB_PASSWORD@localhost:$DB_PORT/$DB_NAME?sslmode=disable"
else
    echo "❌ Failed to connect to PostgreSQL"
    exit 1
fi 