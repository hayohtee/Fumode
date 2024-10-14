CREATE TABLE IF NOT EXISTS cart
(
    cart_id      BIGSERIAL NOT NULL,
    user_id      BIGINT    NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    furniture_id BIGINT    NOT NULL REFERENCES furniture (furniture_id) ON DELETE CASCADE,
    quantity     INTEGER   NOT NULL,
    PRIMARY KEY (cart_id, user_id)
);