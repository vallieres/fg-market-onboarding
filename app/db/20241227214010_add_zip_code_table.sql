-- +goose Up
CREATE TABLE `zipcodes` (
     `zipcode` varchar(10) NOT NULL,
     `city` varchar(128) NOT NULL,
     `state` varchar(128) NOT NULL,
     `created_at` int NOT NULL DEFAULT 0,
     `updated_at` int NOT NULL DEFAULT 0,
     `deleted_at` int,
     KEY `idx_zipcode` (`zipcode`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS zipcodes;
