# Deploying WatchMUD

One Docker host running three containers: the game, its mongo, and a backup job.
Everything below runs from the repo root on that host.

Tested end to end with Docker 29, 2026-09-25: create a character, restart the game
container with that player still connected, log back in in the same room; backups
written and restored; mongo not reachable from the host; a session over TLS, and a
renewed certificate picked up with a player connected, who stayed connected.

## First time

You need Docker with the compose plugin, a checkout of this repo on the host, and a
domain name pointing at the host for the TLS certificate. Get the certificate first
(see TLS, below): the game won't start with TLS configured and no certificate. To
start without TLS, set `tls.port: 0` in `deploy/app.yaml`.

```sh
cp deploy/.env.example deploy/.env
# two long random passwords into deploy/.env:
openssl rand -base64 24 | tr -d '/+='
docker compose -f deploy/compose.yaml up -d --build
docker compose -f deploy/compose.yaml ps        # three services, mongo (healthy)
telnet <host> 4000
```

Open ports 4000 and 4443 in the host's firewall, 80 for certbot's renewals, and
nothing else except ssh. Mongo publishes no port, so it is only reachable from the
other two containers. **Docker writes its own
iptables rules for published ports and goes around `ufw`**: a `ufw deny 4000` won't
close the game port. Use the cloud provider's firewall, or take the port out of
`compose.yaml`.

## TLS

Port 4443 is the same game over TLS, so passwords don't cross the internet in the
clear; port 4000 stays plain telnet for clients that can't do TLS, and tells players
4443 exists. Mudlet has a "Secure" checkbox, TinTin++ has `#ssl`, and anything else
can use `openssl s_client -connect <host>:4443`.

The certificate is Let's Encrypt's, through certbot on the host (not in a container).
certbot's standalone mode answers the challenge on port 80 itself:

```sh
sudo certbot certonly --standalone -d mud.example.com \
  --deploy-hook "$PWD/deploy/certbot-hook.sh"
```

The hook copies the certificate into `deploy/certs`, owned by the game's user (uid
65532 -- `/etc/letsencrypt` is root-only). certbot remembers the hook, and its
renewal timer runs it again every time the certificate is renewed. **The game
reloads the renewed certificate on the next connection, without a restart**, so
renewals don't disconnect anyone. The log says `tls: loaded the renewed certificate`.

Check it from anywhere:

```sh
openssl s_client -connect mud.example.com:4443 </dev/null 2>/dev/null | openssl x509 -noout -dates
```

## Making yourself a wizard

`make wizard` talks to the development mongo, which has no password. Here you go
through the root user. Log out of the game first, or the next timed save writes the
old value back.

```sh
docker compose -f deploy/compose.yaml exec mongo mongosh -u root -p \
  --authenticationDatabase admin watchmud \
  --eval 'db.players.updateOne({name: "Bob"}, {$set: {wizard: true}})'
```

It asks for the password (`MONGO_ROOT_PASSWORD`). The name is capitalized the way the
game stores it: `Bob`, not `bob`.

## Updating

```sh
git pull
docker compose -f deploy/compose.yaml up -d --build watchmud
```

**This disconnects everyone.** The game gets SIGTERM, saves everyone who is logged in,
and exits; compose waits up to 30s (`stop_grace_period`) before it would SIGKILL.
Players reconnect to wherever they were. Announce it first.

## Logs

```sh
docker compose -f deploy/compose.yaml logs -f watchmud
```

Docker keeps them and rotates them: 5 files of 20MB per container.

The game logs where each connection came from (`telnet 1.2.3.4:5678`). After the
first real players connect, **check those are their addresses**. If every
connection comes from the same address (a gateway, a proxy), the 5-per-address
connection cap is counting that address, and the sixth player is refused.

## Backups

The `backup` service writes `deploy/backups/watchmud-<time>.archive.gz` every night
and deletes ones older than two weeks. They're on the same disk as the database, so
they cover a bad deploy or a mistake, not the disk dying. **Copy them off the host**:
an rsync from somewhere else on a cron, or the provider's volume snapshots.

Restoring one replaces the game's data with what's in the file, so stop the game first:

```sh
docker compose -f deploy/compose.yaml stop watchmud
docker compose -f deploy/compose.yaml exec -T mongo mongorestore -u root -p <root password> \
  --authenticationDatabase admin --gzip --archive --drop < deploy/backups/<file>
docker compose -f deploy/compose.yaml start watchmud
```

## Passwords

The game's database password is set when mongo's volume is created, by
`mongo-init.js`, and never again. Changing `WATCHMUD_DB_PASSWORD` in `.env` later
just stops the game connecting. To change it, change it in mongo too:

```sh
docker compose -f deploy/compose.yaml exec mongo mongosh -u root -p \
  --authenticationDatabase admin watchmud \
  --eval 'db.changeUserPassword("watchmud", "<new password>")'
# then set it in .env, and:
docker compose -f deploy/compose.yaml up -d watchmud
```

## Not done yet

- **A health check.** `restart: unless-stopped` brings the game back if it crashes,
  but nothing notices it hanging. A TCP check wouldn't either, and would show up as a
  connection in the logs every time it ran.
