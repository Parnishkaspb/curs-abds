-- +goose Up
-- +goose StatementBegin
CREATE TABLE sources (
    id SERIAL PRIMARY KEY,
    name CHARACTER VARYING(30)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sources;
-- +goose StatementEnd
