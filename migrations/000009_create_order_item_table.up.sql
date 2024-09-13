CREATE TABLE IF NOT EXISTS order_item
(
    order_item_id BIGSERIAL,
    quantity      INT            NOT NULL,
    price         DECIMAL(10, 2) NOT NULL,
    furniture_id  BIGINT REFERENCES furniture (furniture_id),
    order_id      BIGINT         NOT NULL REFERENCES order (order_id) ON DELETE CASCADE
);