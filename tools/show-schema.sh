#!/bin/bash

echo "📊 Database Schema:"
echo ""
echo "=== Tables ==="
docker exec rizon-postgres psql -U rizon -d rizon_db -c "\dt"
echo ""
echo "=== Table Structures ==="
docker exec rizon-postgres psql -U rizon -d rizon_db -t -c "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' AND table_type = 'BASE TABLE' ORDER BY table_name;" | tr -d ' ' | grep -v '^$' | while read table; do
  if [ ! -z "$table" ]; then
    echo ""
    echo "--- Table: $table ---"
    docker exec rizon-postgres psql -U rizon -d rizon_db -c "\d $table"
  fi
done
