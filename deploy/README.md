# Deploying WatchMUD

One Docker host running three containers: the game, its mongo, and a backup job.
Everything below runs from the repo root on that host.

Tested end to end with Docker 29, 2026-09-25: create a character, restart the game
container with that player still connected, log back in in the same room; backups
written and restored; mongo not reachable from the host; a session over TLS, and a
renewed certificate picked up with a player connected, who stayed connected.

## First time

You need Docker with the compose plugin, a checkout of this repo on the host, and a
domain name pointing at the host for the TLS certificate -- production is
`watchmud.com`: telnet on 4000, TLS on 4443. Get the certificate first
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
sudo certbot certonly --standalone -d watchmud.com \
  --deploy-hook "$PWD/deploy/certbot-hook.sh"
```

The hook copies the certificate into `deploy/certs`, owned by the game's user (uid
65532 -- `/etc/letsencrypt` is root-only). certbot remembers the hook, and its
renewal timer runs it again every time the certificate is renewed. **The game
reloads the renewed certificate on the next connection, without a restart**, so
renewals don't disconnect anyone. The log says `tls: loaded the renewed certificate`.

Check it from anywhere:

```sh
openssl s_client -connect watchmud.com:4443 </dev/null 2>/dev/null | openssl x509 -noout -dates
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

The `backup` service writes `deploy/backups/watchmud-<time>.archive.gz` when it starts
and every night after, and deletes ones older than two weeks. Those are on the same
disk as the database: they cover a bad deploy or a mistake, not losing the droplet.

The `offsite` service covers that. Every hour it copies new dumps to a DigitalOcean
Spaces bucket, which keeps them for 90 days. It copies and never syncs, so the droplet
deleting its old dumps doesn't delete them from the bucket.

### Setting up the bucket

1. In the DigitalOcean console, create a Spaces bucket -- say `watchmud-backups` --
   in a **different region from the droplet** (nyc3; the droplet is in sfo3), with
   file listing restricted. The name is global across Spaces, so it may need a suffix.
2. Create a Spaces access key **limited to that bucket**, with read/write/delete
   (delete is how the 90-day pruning works). The secret is shown once.
3. **Before pulling a version of this repo that has the `offsite` service**, add the
   settings to `deploy/.env` -- compose refuses to run anything, even `logs`, while
   they're missing:

   ```sh
   SPACES_REGION=nyc3
   SPACES_BUCKET=watchmud-backups
   SPACES_KEY=<access key>
   SPACES_SECRET=<secret>
   ```

4. Then start it. Only the new service starts; the game isn't touched:

   ```sh
   git pull
   docker compose -f deploy/compose.yaml up -d offsite
   docker compose -f deploy/compose.yaml logs offsite
   ```

   The log should say `Copied (new)` for each dump and `offsite: copied to
   spaces:watchmud-backups/backups`. Check the files are in the bucket in the console.

A key that can delete from the bucket is a key that can empty it, and it's on the
droplet. If the droplet were compromised, the bucket could go too. Turning on the
droplet's own weekly backups in the console is a second copy that key can't touch.

### Restoring

Restoring replaces the game's data with what's in the file, so stop the game first:

```sh
docker compose -f deploy/compose.yaml stop watchmud
docker compose -f deploy/compose.yaml exec -T mongo mongorestore -u root -p <root password> \
  --authenticationDatabase admin --gzip --archive --drop < deploy/backups/<file>
docker compose -f deploy/compose.yaml start watchmud
```

If the droplet is gone, the dump comes from the bucket. Download it in the console, or
on a new droplet with the same `.env`:

```sh
docker compose -f deploy/compose.yaml run --rm -v "$PWD/deploy/backups:/restore" \
  --entrypoint rclone offsite copy spaces:watchmud-backups/backups/<file> /restore
```

then restore it as above.

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
