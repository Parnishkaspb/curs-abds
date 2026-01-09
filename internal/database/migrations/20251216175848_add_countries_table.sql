-- +goose Up
-- +goose StatementBegin
CREATE TABLE countries (
    id SERIAL PRIMARY KEY,
    name CHARACTER VARYING(30)
);

INSERT INTO countries (name) VALUES
    ('RU'),
    ('BY'),
    ('USA'),
    ('GER'),
    ('CHINA');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE countries;
-- +goose StatementEnd

