#!/bin/sh
# Copies the dumps in /backups (deploy/backups on the host) to a DigitalOcean
# Spaces bucket every OFFSITE_EVERY_SECONDS, and deletes ones older than
# OFFSITE_KEEP_DAYS from the bucket. Run by the offsite service in
# deploy/compose.yaml, with the bucket's credentials in its environment as
# the rclone remote "spaces".
#
# copy, not sync: the droplet deleting its two-week-old dumps must not delete
# them here too. The bucket keeps its own, longer, history.
#
# Restore one from the bucket: see deploy/README.md, "Backups".
set -eu

dest="spaces:$SPACES_BUCKET/backups"

while true; do
  if rclone copy /backups "$dest" --include 'watchmud-*.archive.gz' --stats-one-line -v 2>&1; then
    echo "offsite: copied to $dest"
  else
    echo "offsite: FAILED to copy to $dest, will try again next time" >&2
  fi
  rclone delete "$dest" --include 'watchmud-*.archive.gz' --min-age "${OFFSITE_KEEP_DAYS}d" -v 2>&1 \
    || echo "offsite: FAILED to prune old dumps, will try again next time" >&2
  sleep "$OFFSITE_EVERY_SECONDS"
done
