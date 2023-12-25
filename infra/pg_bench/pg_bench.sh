#!/bin/bash

# Impostazioni aggiuntive per pgbench
PG_BENCH_INIT_OPTIONS=""
PG_BENCH_OPTIONS="-c 10 -j 4 -T 60"
export PGPASSWORD=$PG_BEBCHMARK_PASSWORD
echo initdb
echo $PGPASSWORD

# Comando per eseguire pgbench
PG_BENCH_INIT_COMMAND="pgbench -U $PG_BEBCHMARK_USER -h $DB_HOST -p $DB_PORT -s 100 -i -n $PG_BENCH_INIT_OPTIONS $PG_BEBCHMARK_DB"
echo $PG_BENCH_INIT_COMMAND

bash -c "$PG_BENCH_INIT_COMMAND"
touch /tmp/script.end
echo "Initdb completed, waiting 20 seconds waiting for database replication to complete..."
PG_BENCH_COMMAND="pgbench -U $PG_BEBCHMARK_USER -h $DB_HOST -p $DB_PORT -j 20 -r -n $PG_BENCH_OPTIONS $PG_BEBCHMARK_DB"
echo $PG_BENCH_COMMAND
# sleep 200
PG_BENCH_COMMAND="pgbench -U $PG_BEBCHMARK_USER -h $DB_HOST -p $DB_PORT -s 10 -r -n $PG_BENCH_OPTIONS $PG_BEBCHMARK_DB"
echo $PG_BENCH_COMMAND
$PG_BENCH_COMMAND

# Genera 50.000 righe con pgbench
# echo "Generazione di $PG_BEBCHMARK_NUM_ROWS righe con pgbench..."
# sleep 3600
# $PG_BENCH_COMMAND -i
#
# # Esegui il benchmark con pgbench
# echo "Esecuzione del benchmark con pgbench..."
# $PG_BENCH_COMMAND
#
# echo "Benchmark completato. Uscita."
