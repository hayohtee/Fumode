CREATE TABLE IF NOT EXISTS "order"
(
    order_id    BIGSERIAL,
    order_date  TIMESTAMP(0) WITH TIME ZONE,
    total_price DECIMAL(10, 2) NOT NULL,
    customer_id BIGINT         NOT NULL REFERENCES customer (customer_id) ON DELETE CASCADE,
    payment_id  BIGINT         NOT NULL REFERENCES payment (payment_id) ON DELETE CASCADE,
    shipment_id BIGINT         NOT NULL REFERENCES shipment (shipment_id) ON DELETE CASCADE,
    PRIMARY KEY (order_id, customer_id, payment_id, shipment_id)
);