#!/bin/bash

if [ -f .env ]; then
    source .env
fi

go run ./cmd/db_migrations $DB_URL up