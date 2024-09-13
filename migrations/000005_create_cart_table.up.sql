CREATE TABLE IF NOT EXISTS cart
(
    cart_id      BIGSERIAL NOT NULL,
    customer_id  BIGINT    NOT NULL REFERENCES customer (customer_id),
    furniture_id BIGINT    NOT NULL REFERENCES furniture (furniture_id),
    quantity     INTEGER   NOT NULL
);