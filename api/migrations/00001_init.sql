-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp"; -- noqa: L057

CREATE TABLE IF NOT EXISTS users (
    id uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    username varchar(255) NOT NULL UNIQUE,
    password_hash bytea NOT NULL,
    role varchar(20) NOT NULL DEFAULT 'default'
);

CREATE TABLE IF NOT EXISTS tokens (
    hash bytea PRIMARY KEY,
    user_id uuid NOT NULL REFERENCES users ON DELETE CASCADE,
    expiry timestamptz (0) NOT NULL,
    scope int4 NOT NULL,
    used boolean NOT NULL DEFAULT false
);

CREATE TABLE IF NOT EXISTS schools (
    id serial4 PRIMARY KEY,
    name varchar(255) NOT NULL UNIQUE,
    read_only bool NOT NULL DEFAULT false
);
INSERT INTO schools (name, read_only)
VALUES ('Andere', true);

CREATE TABLE IF NOT EXISTS locations (
    id uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    name varchar(255) NOT NULL UNIQUE,
    capacity int4 NOT NULL,
    user_id uuid NOT NULL UNIQUE REFERENCES users ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS check_ins (
    id serial4 PRIMARY KEY,
    location_id uuid NOT NULL REFERENCES locations ON DELETE CASCADE,
    school_id int4 NOT NULL DEFAULT 1 REFERENCES schools ON DELETE SET DEFAULT,
    capacity int4 NOT NULL,
    created_at timestamptz (0) NOT NULL DEFAULT now(),
    created_at_time_zone text CHECK (
        now() AT TIME ZONE created_at_time_zone IS NOT NULL
    )
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS check_ins;
DROP TABLE IF EXISTS locations;
DROP TABLE IF EXISTS schools;
DROP TABLE IF EXISTS tokens;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS roles;
-- +goose StatementEnd
