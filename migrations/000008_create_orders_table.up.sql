CREATE TABLE IF NOT EXISTS orders
(
    order_id    BIGSERIAL PRIMARY KEY,
    order_date  TIMESTAMP(0) WITH TIME ZONE,
    total_price DECIMAL(10, 2) NOT NULL,
    user_id     BIGINT REFERENCES users (user_id),
    payment_id  BIGINT         NOT NULL REFERENCES payment (payment_id) ON DELETE CASCADE,
    shipment_id BIGINT         NOT NULL REFERENCES shipment (shipment_id) ON DELETE CASCADE
);