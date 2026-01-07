-- +goose Up
-- +goose StatementBegin
ALTER TABLE click_logs ADD COLUMN country TEXT;
ALTER TABLE click_logs ADD COLUMN device_type TEXT;
ALTER TABLE click_logs ADD COLUMN traffic TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE click_logs DROP COLUMN country;
ALTER TABLE click_logs DROP COLUMN device_type;
ALTER TABLE click_logs DROP COLUMN traffic;
-- +goose StatementEnd
