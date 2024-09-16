CREATE TABLE IF NOT EXISTS customer
(
    customer_id  BIGSERIAL PRIMARY KEY,
    name         TEXT                        NOT NUlL,
    email        citext UNIQUE               NOT NULL,
    password     BYTEA,
    address      TEXT,
    phone_number VARCHAR(100),
    role         VARCHAR(15)                 NOT NULL DEFAULT 'customer',
    created_at   TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);