-- +goose Up
-- +goose StatementBegin
INSERT INTO scraping_key (name, created_at, updated_at) VALUES 
('important-key', NOW(), NOW()),
('no-this-important-key', NOW(), NOW()),
('no-this-more-important-key', NOW(), NOW());
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM scraping_key WHERE name IN ('important-key', 'no-this-important-key', 'no-this-more-important-key');
-- +goose StatementEnd
