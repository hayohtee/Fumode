CREATE TABLE IF NOT EXISTS shipment
(
    shipment_id   BIGSERIAL,
    shipment_date TIMESTAMP(0) WITH TIME ZONE,
    address       TEXT         NOT NULL,
    city          VARCHAR(100) NOT NULL,
    state         VARCHAR(100) NOT NULL,
    country       VARCHAR(100) NOT NULL,
    zip_code      VARCHAR(10)  NOT NULL,
    customer_id   BIGINT       NOT NULL REFERENCES customer (customer_id) ON DELETE CASCADE,
    PRIMARY KEY (shipment_id, customer_id)
);