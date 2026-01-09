-- +goose Up
-- +goose StatementBegin
CREATE TABLE currencies (
    id SERIAL PRIMARY KEY,
    name CHARACTER VARYING(100),
    iso CHARACTER VARYING(50)
);

INSERT INTO currencies (name, iso) VALUES
                                       ('RU', 'RU'),
                                       ('BY', 'BY'),
                                       ('USA', 'USD'),
                                       ('EUR', 'GER'),
                                       ('YEN', 'CHINA');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE currencies;
-- +goose StatementEnd
