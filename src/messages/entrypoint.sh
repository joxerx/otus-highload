#!/bin/sh

echo "Waiting for postgres..."

while ! nc -z $POSTGRES_HOST $POSTGRES_MASTER_PORT; do
  sleep 0.1
done

echo "PostgreSQL started"

exec "$@"