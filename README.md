# Postgresql Health Check

This is a simple server that can be used in combination with HAproxy to check which database is the primary and which is a follower.

The main goal is to replace xinetd scripts with something more structured and ready to use, without having to copy and paste the same bash scripts.

## Installation

Depending on the distribution, you will need to choose the appropriate RPM or DEB package, taking into account your architecture. Typically, x86_64 is the most common choice. For RPM-based systems, you can download the RPM package from the release and install it with the command:

```bash
rpm -ivh postgresql-health-check-0.1.0-1.x86_64.rpm
```

For other distributions, you can download the tarball and install it manually.

## Configuration

The configuration is straightforward and is located in the file `/etc/postgresql-health-check/config.yml`.

```yaml
---
postgres:
  host: localhost
  port: 5432
  user: pgsql_health_check_monitor
  password: veryS3cr3t
  dbname: postgres
server:
  port: 8080
  host: localhost
  mode: release
  # - debug
  # - release
```

## Usage
To collect metrics from the `pg_stat*` views as a non-superuser in PostgreSQL server versions >= 10, you can assign the roles pg_monitor or pg_read_all_stats to the user.

Example:

```sql
CREATE USER pgsql_health_check_monitor WITH PASSWORD 'veryS3cr3t';
GRANT pg_monitor TO pgsql_health_check_monitor;
```

### Check if the DB is in write mode

```bash
curl http://localhost:8080/write
```

### Check if the DB is in read mode

```bash
curl http://localhost:8080/read
```

## HAproxy Config

pg-health-check is designed to be used with HAproxy, although I am a big fan of nginx, HAproxy is more suitable for this type of task because it supports active checks, unlike nginx, which only supports passive checks in the community version.
If you have nginx-plus, you can use active checks, but that's not covered at now.

If you have nginx-plus, open a PR and I will be happy to merge it with new documentation.

```bash
global
    maxconn 200

defaults
    log stdout format raw local0
    mode tcp
    retries 2
    timeout client 10m
    timeout connect 5s
    timeout server 10m
    timeout check 5s

listen stats
    mode http
    bind *:8085
    stats enable
    stats uri /stats

listen Write
    bind *:5435
    option httpchk
    http-check send meth GET  uri /write
    http-check expect status 200
    default-server inter 3s fall 3 rise 2 on-marked-down shutdown-sessions
    server pg-data-1 pg-data-1:5432 maxconn 100 check port 8080
    server pg-data-2 pg-data-2:5432 maxconn 100 check port 8080

listen Read
    bind *:5436
    option httpchk
    http-check send meth GET  uri /read
    http-check expect status 200
    default-server inter 3s fall 3 rise 2 on-marked-down shutdown-sessions
    server pg-data-1 pg-data-1:5432 maxconn 100 check port 8080
    server pg-data-2 pg-data-2:5432 maxconn 100 check port 8080
```