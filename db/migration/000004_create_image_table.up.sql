CREATE TABLE image
(
    id      BIGSERIAL PRIMARY KEY,
    url     VARCHAR(512) NOT NULL UNIQUE,
    spot_id BIGINT       NOT NULL,
    CONSTRAINT fk_spot_id FOREIGN KEY (spot_id) REFERENCES spot (id)
);
