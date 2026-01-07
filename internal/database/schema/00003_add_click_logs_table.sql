-- +goose Up
-- +goose StatementBegin
CREATE TABLE click_logs (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    link_id     UUID NOT NULL,
    ip_address  VARCHAR(255),
    user_agent  VARCHAR(255),
    referrer    VARCHAR(255),
    clicked_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    FOREIGN KEY (link_id) REFERENCES links(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE click_logs;
-- +goose StatementEnd
