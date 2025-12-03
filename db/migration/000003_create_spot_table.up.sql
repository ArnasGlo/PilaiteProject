CREATE TYPE spot_category AS ENUM ('Gamta', 'Lauko_treniruokliai', 'Slaptos_vietos', 'Restoranai', 'Parduotuves');

CREATE TABLE spot
(
    id            BIGSERIAL PRIMARY KEY,
    category      spot_category NOT NULL ,
    name          VARCHAR(50)  NOT NULL,
    description   VARCHAR(255) NOT NULL,
    location_id   BIGINT       NOT NULL UNIQUE,
    CONSTRAINT fk_location_id FOREIGN KEY (location_id) REFERENCES location (id)
);
