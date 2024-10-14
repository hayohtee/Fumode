CREATE TABLE IF NOT EXISTS review
(
    review_id    BIGSERIAL PRIMARY KEY,
    user_id      BIGINT                      NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    furniture_id BIGINT                      NOT NULL REFERENCES furniture (furniture_id) ON DELETE CASCADE,
    rating       INTEGER                     NOT NULL,
    comment      TEXT                        NOT NULL,
    created_at   TIMESTAMP(0) WITH TIME ZONE NOT NULL,
    version      INTEGER                     NOT NULL DEFAULT 1
);