#!/bin/bash

# Load environment variables from .env file
if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
fi

# Check if required env vars are set
if [ -z "$DB_HOST" ] || [ -z "$DB_USER" ] || [ -z "$DB_NAME" ]; then
    echo "Error: Required environment variables not set"
    echo "Please ensure DB_HOST, DB_USER, DB_NAME are defined in .env"
    exit 1
fi

echo "======================================"
echo "Running all migrations"
echo "Database: $DB_NAME on $DB_HOST"
echo "======================================"

# Get all migration files and sort them
MIGRATION_FILES=$(ls -1 migrations/*.sql 2>/dev/null | sort)

if [ -z "$MIGRATION_FILES" ]; then
    echo "Error: No migration files found in migrations/ directory"
    exit 1
fi

# Run each migration in order
for file in $MIGRATION_FILES; do
    echo ""
    echo "Running: $file"
    echo "--------------------------------------"
    
    PGPASSWORD="$DB_PASSWORD" psql \
        -h "$DB_HOST" \
        -p "${DB_PORT:-5432}" \
        -U "$DB_USER" \
        -d "$DB_NAME" \
        -f "$file"
    
    if [ $? -ne 0 ]; then
        echo ""
        echo "ERROR: Migration failed: $file"
        echo "======================================"
        exit 1
    fi
    
    echo "✓ Completed: $file"
done

echo ""
echo "======================================"
echo "All migrations completed successfully!"
echo "======================================"
