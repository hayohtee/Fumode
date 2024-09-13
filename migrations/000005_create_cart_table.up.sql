CREATE TABLE IF NOT EXISTS cart
(
    cart_id      BIGSERIAL NOT NULL,
    customer_id  BIGINT    NOT NULL REFERENCES customer (customer_id) ON DELETE CASCADE,
    furniture_id BIGINT    NOT NULL REFERENCES furniture (furniture_id) ON DELETE CASCADE,
    quantity     INTEGER   NOT NULL,
    PRIMARY KEY (cart_id, customer_id, furniture_id)
);