CREATE TABLE location
(
    id        BIGSERIAL PRIMARY KEY,
    address   VARCHAR(30)    NOT NULL,
    latitude  DECIMAL(10, 8) NOT NULL,
    longitude DECIMAL(11, 8) NOT NULL
);
