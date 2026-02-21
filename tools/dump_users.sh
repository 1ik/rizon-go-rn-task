#!/bin/bash
# Dump users table: last 5 records, descending order

echo "📋 Users (last 5, newest first):"
echo ""
docker exec rizon-postgres psql -U rizon -d rizon_db -c "SELECT * FROM users ORDER BY id DESC LIMIT 5;"
