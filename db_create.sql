DROP DATABASE IF EXISTS intelliada_db;
DROP ROLE IF EXISTS intelliada_db;

CREATE DATABASE intelliada_db;
CREATE USER admin WITH password '';
GRANT ALL privileges ON DATABASE fias TO dev;
