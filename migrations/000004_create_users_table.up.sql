CREATE TABLE IF NOT EXISTS users (
    id bigserial PAIMARY KEY,
    created_at timestamptz NOT NULL DEFAULT now(),
    name text NOT NULL,
    email citext NOT NULL,
    password_has bytea NOT NULL,
    activated bool NOT NULL,
    version integer NOT NULL DEFAULT 1
);

