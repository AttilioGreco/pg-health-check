#!/bin/bash

su - postgres

whoami

echo "Creating .pgpass file for replication user"
echo "pg-data-1:5432:replication:repuser:replication_password" > $HOME/.pgpass
chmod 0600 $HOME/.pgpass
pg_basebackup -d 'host=pg-data-1 user=repuser' -D /var/lib/postgresql/data/ -R -P

# echo "Moving .pgpass file to /var/lib/postgresql"
# mv $HOME/.pgpass /var/lib/postgresql/.pgpass

echo "Changing ownership of /var/lib/postgresql to postgres"
chown -R postgres:postgres /var/lib/postgresql

echo "Start PostgreSQL"
/usr/local/bin/docker-entrypoint.sh postgres
#
# sleep 3000
#
# echo "Base backup completed startup process"