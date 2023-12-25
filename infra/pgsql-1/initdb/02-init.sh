#!/bin/bash

ip_subnet=$(ip a | grep inet | grep -v 127.0.0.1 | awk '{print $2}')
echo "host    replication     repuser     $ip_subnet   scram-sha-256" >> /var/lib/postgresql/data/pg_hba.conf
echo "Reload postgresql to apply new pg_hba.conf"
pg_ctl reload