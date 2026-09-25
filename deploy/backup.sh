#!/usr/bin/env bash
# Dumps the watchmud database into /backups (deploy/backups on the host) every
# BACKUP_EVERY_SECONDS, and deletes dumps older than BACKUP_KEEP_DAYS.
# Run by the backup service in deploy/compose.yaml.
#
# Restore one with:
#   docker compose -f deploy/compose.yaml exec -T mongo \
#     mongorestore --username root --password "$MONGO_ROOT_PASSWORD" \
#     --authenticationDatabase admin --gzip --archive --drop < deploy/backups/<file>
set -euo pipefail

# The password goes in a config file rather than on the command line, where
# anything that can list processes could read it.
conf=$(mktemp)
trap 'rm -f "$conf"' EXIT
printf 'password: %s\n' "$MONGO_ROOT_PASSWORD" > "$conf"

while true; do
  file="/backups/watchmud-$(date -u +%Y%m%dT%H%M%SZ).archive.gz"
  if mongodump --host mongo --username root --config "$conf" \
       --authenticationDatabase admin --db watchmud \
       --gzip --archive="$file" --quiet; then
    echo "backup: wrote $file ($(du -h "$file" | cut -f1))"
  else
    echo "backup: FAILED, will try again next time" >&2
    rm -f "$file"
  fi
  find /backups -name 'watchmud-*.archive.gz' -mtime +"$BACKUP_KEEP_DAYS" -print -delete
  sleep "$BACKUP_EVERY_SECONDS"
done
