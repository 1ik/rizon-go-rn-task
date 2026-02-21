#!/bin/bash
# Dump feedbacks table: last 5 records, descending order

echo "📋 Feedbacks (last 5, newest first):"
echo ""
docker exec rizon-postgres psql -U rizon -d rizon_db -c "SELECT * FROM feedbacks ORDER BY id DESC LIMIT 5;"
