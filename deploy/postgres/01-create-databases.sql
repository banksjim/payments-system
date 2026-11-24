-- Runs once, when the Postgres data volume is first created.
-- One database per service: each service gets its own database.
CREATE DATABASE wallet;
CREATE DATABASE auth;
CREATE DATABASE transfers;
