# Deploying WatchMUD

One Docker host running three containers: the game, its mongo, and a backup job.
Everything below runs from the repo root on that host.

Tested end to end with Docker 29, 2026-09-25: create a character, restart the game
container with that player still connected, log back in in the same room; backups
written and restored; mongo not reachable from the host.

## First time

You need Docker with the compose plugin, and a checkout of this repo on the host.

```sh
cp deploy/.env.example deploy/.env
# two long random passwords into deploy/.env:
openssl rand -base64 24 | tr -d '/+='
docker compose -f deploy/compose.yaml up -d --build
docker compose -f deploy/compose.yaml ps        # three services, mongo (healthy)
telnet <host> 4000
```

Open port 4000 in the host's firewall, and nothing else except ssh. Mongo publishes no
port, so it is only reachable from the other two containers. **Docker writes its own
iptables rules for published ports and goes around `ufw`**: a `ufw deny 4000` won't
close the game port. Use the cloud provider's firewall, or take the port out of
`compose.yaml`.

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

- **TLS.** Passwords still cross the internet as plain telnet. See ROADMAP, Deploy.
- **A health check.** `restart: unless-stopped` brings the game back if it crashes,
  but nothing notices it hanging. A TCP check wouldn't either, and would show up as a
  connection in the logs every time it ran.
