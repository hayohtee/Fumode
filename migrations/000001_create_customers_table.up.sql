CREATE TABLE IF NOT EXISTS customers
(
    customer_id  BIGSERIAL PRIMARY KEY,
    name         TEXT          NOT NUlL,
    email        citext UNIQUE NOT NULL,
    password     BYTEA,
    address      TEXT,
    phone_number VARCHAR(100),
    created_at   TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    version      int           NOT NULL DEFAULT 1
);