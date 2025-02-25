-- +goose Up
-- +goose StatementBegin
CREATE INDEX container_name_idx ON containers (name)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index if exists container_name_idx;
-- +goose StatementEnd
