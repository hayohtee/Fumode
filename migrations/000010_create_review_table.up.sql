CREATE TABLE IF NOT EXISTS review
(
    review_id    BIGSERIAL PRIMARY KEY,
    customer_id  BIGINT                      NOT NULL REFERENCES customer (customer_id),
    furniture_id BIGINT                      NOT NULL REFERENCES furniture (furniture_id),
    rating       INTEGER                     NOT NULL,
    comment      TEXT                        NOT NULL,
    created_at   TIMESTAMP(0) WITH TIME ZONE NOT NULL,
    version      INTEGER                     NOT NULL DEFAULT 1
);