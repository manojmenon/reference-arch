#!/bin/sh
# One-way logical replication: postgres-primary-1 is the write leader (matches the API
# DATABASE_URL). postgres-primary-2 and postgres-primary-3 subscribe only to node 1.
#
# A full mesh (each node subscribing to every other) causes duplicate INSERTs on the same
# row: e.g. an insert on node 2 reaches node 1 via 2→1 and again via 2→3→1, violating PKs.
# Native PostgreSQL logical replication is not a turnkey multi-master solution for the same
# table; use a single publication + fan-out, or a product like BDR for true multi-master.
set -e

REPL_PW="${REPL_PASSWORD:-replpass}"
export PGPASSWORD="${POSTGRES_PASSWORD:-apppass}"

wait_for() {
  host="$1"
  until pg_isready -h "$host" -p 5432 -U appuser -d appdb -q; do
    echo "waiting for $host..."
    sleep 2
  done
}

for h in postgres-primary-1 postgres-primary-2 postgres-primary-3; do
  wait_for "$h"
done

# Publication only on the write primary (backend connects here).
psql -h postgres-primary-1 -U appuser -d appdb -v ON_ERROR_STOP=1 <<'EOSQL'
DO $do$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_catalog.pg_publication WHERE pubname = 'pub_users') THEN
    EXECUTE 'CREATE PUBLICATION pub_users FOR TABLE users';
  END IF;
END
$do$;
EOSQL

run_sub() {
  target="$1"
  name="$2"
  source="$3"
  sql="CREATE SUBSCRIPTION ${name} CONNECTION 'host=${source} port=5432 user=replicator password=${REPL_PW} dbname=appdb' PUBLICATION pub_users WITH (copy_data = off);"
  psql -h "$target" -U appuser -d appdb -c "$sql" || true
}

# Unique subscription names on publisher postgres-primary-1 (one slot per subscriber).
run_sub postgres-primary-2 sub_pg2_from_pg1 postgres-primary-1
run_sub postgres-primary-3 sub_pg3_from_pg1 postgres-primary-1

echo "Replication bootstrap finished (fan-out from postgres-primary-1)."
