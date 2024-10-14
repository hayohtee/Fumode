CREATE TABLE IF NOT EXISTS wishlist
(
    wishlist_id  BIGSERIAL,
    user_id  BIGINT NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    furniture_id BIGINT NOT NULL REFERENCES furniture (furniture_id) ON DELETE CASCADE,
    PRIMARY KEY (wishlist_id, user_id)
);