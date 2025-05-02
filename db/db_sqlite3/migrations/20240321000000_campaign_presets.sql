-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS "campaign_presets" (
    "id" INTEGER PRIMARY KEY AUTOINCREMENT,
    "user_id" INTEGER NOT NULL,
    "name" TEXT NOT NULL,
    "created_date" DATETIME NOT NULL,
    "template_id" INTEGER,
    "page_id" INTEGER,
    "smtp_id" INTEGER,
    "url" TEXT
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE "campaign_presets"; 