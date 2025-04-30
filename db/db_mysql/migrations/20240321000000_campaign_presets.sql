-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS `campaign_presets` (
    `id` BIGINT(20) NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT(20) NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `created_date` DATETIME NOT NULL,
    `template_id` BIGINT(20),
    `page_id` BIGINT(20),
    `smtp_id` BIGINT(20),
    `url` TEXT,
    PRIMARY KEY (`id`)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `campaign_presets`; 