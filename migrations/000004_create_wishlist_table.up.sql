CREATE TABLE IF NOT EXISTS wishlist
(
    wishlist_id  BIGSERIAL,
    customer_id  BIGINT NOT NULL REFERENCES customer (customer_id) ON DELETE CASCADE,
    furniture_id BIGINT NOT NULL REFERENCES furniture (furniture_id) ON DELETE CASCADE,
    PRIMARY KEY (wishlist_id, customer_id, furniture_id)
);