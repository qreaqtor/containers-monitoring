-- +goose Up
-- +goose StatementBegin
CREATE INDEX container_updated_idx ON containers (updated_at)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index if exists container_updated_idx;
-- +goose StatementEnd
