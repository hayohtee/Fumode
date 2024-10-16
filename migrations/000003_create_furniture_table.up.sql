CREATE TABLE IF NOT EXISTS furniture
(
    furniture_id BIGSERIAL PRIMARY KEY,
    name         VARCHAR(100)   NOT NULL,
    description  TEXT           NOT NULL,
    price        DECIMAL(10, 2) NOT NULL,
    stock        INTEGER        NOT NULL,
    banner_url   TEXT           NOT NULL,
    image_urls   TEXT[],
    category_id  BIGINT REFERENCES category (category_id) ON UPDATE CASCADE,
    version      INTEGER        NOT NULL DEFAULT 1
);