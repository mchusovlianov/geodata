-- Version: 1.1
-- Description: Create table countries
CREATE TABLE countries
(
    uuid         VARCHAR(36),
    code         VARCHAR(2),
    name         VARCHAR(56),
    date_created DATETIME,
    date_updated DATETIME,

    PRIMARY KEY (uuid)
);

-- Version: 1.2
-- Description: Add unique index for name  field
ALTER TABLE countries
    ADD UNIQUE INDEX index_country_name (name);

-- Version: 1.3
-- Description: Add unique index for code field
ALTER TABLE countries
    ADD UNIQUE INDEX index_country_code (code);


-- Version: 1.4
-- Description: Create table cities
CREATE TABLE cities
(
    uuid         VARCHAR(36),
    country_uuid VARCHAR(36),
    name         VARCHAR(176),
    date_created DATETIME,
    date_updated DATETIME,

    PRIMARY KEY (uuid)
);

-- Version: 1.5
-- Description: Add index for country_uuid field
ALTER TABLE cities
    ADD INDEX index_country_uuid (country_uuid);

-- Version: 1.6
-- Description: Add unique index for name and country_uuid field
ALTER TABLE cities
    ADD UNIQUE INDEX index_country_uuid_and_name (country_uuid, name);

-- Version: 1.7
-- Description: Create table locations
CREATE TABLE locations
(
    uuid          VARCHAR(36),
    city_uuid     VARCHAR(36),
    ip            VARCHAR(39), -- max size for IPv4 and IPv6 versions
    mystery_value BIGINT,
    latitude      DOUBLE,
    longitude     DOUBLE,
    date_created  DATETIME,
    date_updated  DATETIME,

    PRIMARY KEY (uuid)
);

-- Version: 1.8
-- Description: Add index for country_uuid field
ALTER TABLE locations
    ADD INDEX index_city_uuid (city_uuid);

-- Version: 1.9
-- Description: Add index for country_uuid field
ALTER TABLE locations
    ADD UNIQUE INDEX index_city_uuid_ip_latitude_longitude (city_uuid, ip, latitude, longitude);

-- Version: 2
-- Description: Add index for location_ip field
ALTER TABLE locations
    ADD UNIQUE INDEX index_ip (ip);
