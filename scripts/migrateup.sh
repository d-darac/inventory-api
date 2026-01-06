#!/bin/bash

if [ -f .env ]; then
    source .env
fi

go run ./cmd/db_migrations up