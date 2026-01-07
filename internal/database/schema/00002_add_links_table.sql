-- +goose Up
-- +goose StatementBegin
CREATE TABLE links (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    original_url        VARCHAR(255) NOT NULL,
    short_code          VARCHAR(255) NOT NULL UNIQUE,
    custom_short_code   VARCHAR(255) UNIQUE,
    user_id             UUID NOT NULL,
    expired_at          TIMESTAMPTZ,

    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at          TIMESTAMPTZ,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE cascade ON UPDATE cascade
);

CREATE INDEX idx_short_code ON links(short_code);
CREATE INDEX idx_custom_short_code ON links(custom_short_code);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE links;
-- +goose StatementEnd
