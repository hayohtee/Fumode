CREATE TABLE IF NOT EXISTS payment
(
    payment_id     BIGSERIAL,
    payment_date   TIMESTAMP(0) WITH TIME ZONE,
    payment_method VARCHAR(100),
    amount         DECIMAL(10, 2),
    customer_id    BIGINT NOT NULL REFERENCES customer (customer_id),
    PRIMARY KEY (payment_id, customer_id)
);