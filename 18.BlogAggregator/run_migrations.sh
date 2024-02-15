#! /bin/sh
# run_migrations.sh
goose -dir sql/schema postgres "${CONN}" "$1"
