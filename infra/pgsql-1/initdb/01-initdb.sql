CREATE USER benckmark with
NOSUPERUSER
CREATEDB
NOCREATEROLE
LOGIN
encrypted password 'benckmark';

CREATE DATABASE benckmark;
ALTER DATABASE benckmark OWNER TO benckmark;

CREATE USER repuser
    REPLICATION
    LOGIN CONNECTION LIMIT 10
ENCRYPTED PASSWORD 'replication_password';