-- +migrate Up
CREATE TABLE IF NOT EXISTS admins (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMPTZ NULL
);

-- +migrate Down
DROP TABLE IF EXISTS admins;