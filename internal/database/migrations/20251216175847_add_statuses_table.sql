-- +goose Up
-- +goose StatementBegin
CREATE TABLE statuses (
    id SERIAL PRIMARY KEY,
    name CHARACTER VARYING(30)
);

INSERT INTO statuses (name) VALUES
                                 ('APPROVED'),
                                 ('DECLINED');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE statuses;
-- +goose StatementEnd

