-- +goose Up
-- +goose StatementBegin
ALTER TABLE click_logs ADD COLUMN browser TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE click_logs DROP COLUMN browser;
-- +goose StatementEnd
