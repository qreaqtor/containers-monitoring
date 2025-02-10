-- +goose Up
-- +goose StatementBegin
CREATE TABLE containers (
    name TEXT PRIMARY KEY,
    id TEXT NOT NULL,
    image TEXT NOT NULL,
    ip TEXT,
    ports TEXT[],
    state TEXT NOT NULL,
    status TEXT NOT NULL,
    updated_at TIMESTAMP DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists containers;
-- +goose StatementEnd
