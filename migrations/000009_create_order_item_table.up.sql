CREATE TABLE IF NOT EXISTS order_item
(
    order_item_id     BIGSERIAL,
    quantity          INT            NOT NULL,
    price             DECIMAL(10, 2) NOT NULL,
    furniture_id      BIGINT REFERENCES furniture (furniture_id),
    customer_order_id BIGINT         NOT NULL REFERENCES customer_order (customer_order_id) ON DELETE CASCADE,
    PRIMARY KEY (order_item_id, order_id)
);