CREATE DATABASE filestask;
CREATE TABLE files (
    id SERIAL PRIMARY KEY,
    filename VARCHAR(80) NOT NULL,
    path VARCHAR(80) NOT NULL,
    modtime DATE,
    checksum BYTEA,
    size INTEGER,
    Uuid uuid
);