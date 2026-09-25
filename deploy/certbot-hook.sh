#!/bin/sh
# certbot's --deploy-hook: copies a new or renewed certificate into
# deploy/certs, where the game container reads it. The game notices the
# changed files on its next TLS handshake; no restart, nobody disconnected.
#
# It exists because /etc/letsencrypt is readable by root only, and the game
# runs as distroless's nonroot user, uid 65532. See deploy/README.md, "TLS".
set -eu
certs="$(cd "$(dirname "$0")" && pwd)/certs"
mkdir -p "$certs"
# key first: the game only reloads once both files have changed and load as
# a pair, so a handshake between the two copies keeps the old certificate
install -m 0600 -o 65532 -g 65532 "$RENEWED_LINEAGE/privkey.pem" "$certs/privkey.pem"
install -m 0644 -o 65532 -g 65532 "$RENEWED_LINEAGE/fullchain.pem" "$certs/fullchain.pem"
echo "watchmud: certificate installed in $certs"
