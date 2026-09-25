// Runs once, the first time mongo starts with an empty volume (the image runs
// everything in /docker-entrypoint-initdb.d then, and never again). It makes
// the user the game connects as: read and write on its own database, nothing
// else -- the root user from MONGO_INITDB_ROOT_* is for you, not the game.
//
// Changing WATCHMUD_DB_PASSWORD in .env later does NOT change this user's
// password; the script doesn't run again. See deploy/README.md.
db.getSiblingDB("watchmud").createUser({
  user: "watchmud",
  pwd: process.env.WATCHMUD_DB_PASSWORD,
  roles: [{ role: "readWrite", db: "watchmud" }],
});
