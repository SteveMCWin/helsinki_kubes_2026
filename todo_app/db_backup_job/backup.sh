#!/usr/bin/env bash
set -euo pipefail

# build a timestamped filename for this backup
TIMESTAMP=$(date -u +%Y%m%dT%H%M%SZ)
FILENAME="todo-db-${TIMESTAMP}.sql.gz"

# dump the database and compress it straight to /tmp
PGPASSWORD="$DB_PASS" pg_dump -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" \
  | gzip > "/tmp/${FILENAME}"

# gcloud doesn't read GOOGLE_APPLICATION_CREDENTIALS on its own, so log in explicitly
gcloud auth activate-service-account --key-file="$GOOGLE_APPLICATION_CREDENTIALS"

# upload the backup to the bucket
gcloud storage cp "/tmp/${FILENAME}" "${GCS_DEST%/}/${FILENAME}"
