-- +goose Up
-- +goose StatementBegin
CREATE TABLE `plans` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(128) NOT NULL,
    `status` VARCHAR(16) NOT NULL DEFAULT 'NEW',
    `pet_name` VARCHAR(64) NOT NULL,
    `pet_species` VARCHAR(64),
    `pet_breed` VARCHAR(64),
    `pet_weight_lbs` VARCHAR(128),
    `pet_activity_level` VARCHAR(128),
    `daily_total_calories` FLOAT,
    `daily_total_protein` FLOAT,
    `daily_total_fat` FLOAT,
    `daily_total_carbohydrates` FLOAT,
    `protein_percent_of_calories` FLOAT,
    `fat_percent_of_calories` FLOAT,
    `carbohydrates_percent_of_calories` FLOAT,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE `weekly_totals` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `plan_id` BIGINT NOT NULL,
    `product` VARCHAR(128) NOT NULL,
    `product_variant_id` VARCHAR(128) NOT NULL,
    `daily_weight_lbs` FLOAT,
    `weekly_weight_lbs` FLOAT,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`plan_id`) REFERENCES `plans`(`id`),
    KEY `idx_deleted_at` (`deleted_at`),
    KEY `idx_plan_id` (`plan_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE `requirements` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `plan_id` BIGINT NOT NULL,
    `description` VARCHAR(128) NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`plan_id`) REFERENCES `plans`(`id`),
    KEY `idx_deleted_at` (`deleted_at`),
    KEY `idx_plan_id` (`plan_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE `recommendations` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `plan_id` BIGINT NOT NULL,
    `description` VARCHAR(128) NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`plan_id`) REFERENCES `plans`(`id`),
    KEY `idx_deleted_at` (`deleted_at`),
    KEY `idx_plan_id` (`plan_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE `meals` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `plan_id` BIGINT NOT NULL,
    `name` VARCHAR(128) NOT NULL,
    `product` VARCHAR(128) NOT NULL,
    `product_variant_id` VARCHAR(128) NOT NULL,
    `quantity_grams` INT NOT NULL,
    `quantity_oz` FLOAT NOT NULL,
    `protein_grams` FLOAT NOT NULL,
    `fat_grams` FLOAT NOT NULL,
    `carbohydrates_grams` FLOAT NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`plan_id`) REFERENCES `plans`(`id`),
    KEY `idx_deleted_at` (`deleted_at`),
    KEY `idx_plan_id` (`plan_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS meals;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS recommendations;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS requirements;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS weekly_totals;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS plans;
-- +goose StatementEnd
