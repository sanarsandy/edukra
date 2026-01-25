#!/bin/bash

# Script Migrasi Khusus untuk Edukra LMS (Docker Environment)
# Configuration based on user provided docker-compose
CONTAINER_NAME="edukra-db-prod"
DB_USER="edukra_lms"
DB_NAME="edukra_db"
# Password default jika tidak diset di env var
export PGPASSWORD="${DB_PASSWORD:-Edukr4LMS}"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Starting Edukra Migration Process...${NC}"
echo "Target Container: $CONTAINER_NAME"
echo "Database: $DB_NAME (User: $DB_USER)"

# 1. Check if container is running
if ! docker ps --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
    echo -e "${RED}Error: Container '$CONTAINER_NAME' is not running!${NC}"
    echo "Please start your docker services first."
    exit 1
fi

# 2. Locate Migrations Directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
MIGRATIONS_DIR="$SCRIPT_DIR/../migrations"

if [ ! -d "$MIGRATIONS_DIR" ]; then
    echo -e "${RED}Error: Migrations directory not found at $MIGRATIONS_DIR${NC}"
    exit 1
fi

# 3. Find and Sort Migration Files
# Using sort -V to handle version numbers correctly (e.g. 2 vs 10)
MIGRATIONS=$(find "$MIGRATIONS_DIR" -name "*.sql" -maxdepth 1 | sort -V)

if [ -z "$MIGRATIONS" ]; then
    echo -e "${RED}No migration files found in $MIGRATIONS_DIR${NC}"
    exit 0
fi

# 4. Execute Migrations
count=0
success=0
skipped=0

for migration_file in $MIGRATIONS; do
    filename=$(basename "$migration_file")
    ((count++))
    
    echo -n "Processing [$count] $filename ... "
    
    # Execute SQL inside the container
    # We strip comments to avoid issues, though psql usually handles them fine
    OUTPUT=$(docker exec -i "$CONTAINER_NAME" psql -U "$DB_USER" -d "$DB_NAME" < "$migration_file" 2>&1)
    EXIT_CODE=$?

    if [ $EXIT_CODE -eq 0 ]; then
        echo -e "${GREEN}DONE${NC}"
        ((success++))
    else
        # Check for common idempotency errors (already exists)
        if echo "$OUTPUT" | grep -q -E "already exists|duplicate|relation .* exists"; then
            echo -e "${YELLOW}SKIPPED (Already applied)${NC}"
            ((skipped++))
        else
            echo -e "${RED}FAILED${NC}"
            echo -e "Error output:\n$OUTPUT"
            # Optional: Exit on first error?
             exit 1
        fi
    fi
done

echo "----------------------------------------"
echo -e "${GREEN}Migration Summary:${NC}"
echo "Total Files Scanned: $count"
echo "Successfully Executed: $success"
echo "Skipped (Existing): $skipped"

if [ $success -eq 0 ] && [ $skipped -eq 0 ]; then
     echo -e "${RED}Warning: No migrations were applied/checked.${NC}"
else
     echo -e "${GREEN}All checks passed.${NC}"
fi
