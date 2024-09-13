CREATE TABLE IF NOT EXISTS categories
(
    category_id BIGSERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    version     INTEGER      NOT NULL DEFAULT 1
);