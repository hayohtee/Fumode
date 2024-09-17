CREATE TABLE IF NOT EXISTS admin
(
    admin_id   BIGSERIAL PRIMARY KEY,
    name       TEXT                        NOT NULL,
    email      citext UNIQUE               NOT NULL,
    password   bytea                       NOT NULL,
    role       VARCHAR(15)                 NOT NULL DEFAULT 'admin',
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);