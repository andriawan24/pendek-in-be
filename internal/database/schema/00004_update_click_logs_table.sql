-- +goose Up
-- +goose StatementBegin
ALTER TABLE click_logs DROP COLUMN link_id;
ALTER TABLE click_logs ADD COLUMN code VARCHAR(255) NOT NULL;

CREATE OR REPLACE FUNCTION check_code_exists(p_code VARCHAR)
RETURNS BOOLEAN AS $$
BEGIN
    RETURN EXISTS (
        SELECT 1 FROM links 
        WHERE short_code = p_code OR custom_short_code = p_code
    );
END;
$$ LANGUAGE plpgsql;

ALTER TABLE click_logs ADD CONSTRAINT chk_code_exists 
    CHECK (check_code_exists(code));

CREATE INDEX idx_click_logs_code ON click_logs(code);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_click_logs_code;
ALTER TABLE click_logs DROP CONSTRAINT IF EXISTS chk_code_exists;
DROP FUNCTION IF EXISTS check_code_exists(VARCHAR);
ALTER TABLE click_logs DROP COLUMN IF EXISTS code;
ALTER TABLE click_logs ADD COLUMN link_id UUID;
-- +goose StatementEnd
