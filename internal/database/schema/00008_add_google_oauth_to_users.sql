-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN google_id VARCHAR(255) UNIQUE;
CREATE INDEX idx_users_google_id ON users (google_id) WHERE google_id IS NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_users_google_id;
ALTER TABLE users DROP COLUMN IF EXISTS google_id;
-- +goose StatementEnd
