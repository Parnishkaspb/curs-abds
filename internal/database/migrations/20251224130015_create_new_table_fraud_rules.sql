-- +goose Up
-- +goose StatementBegin
CREATE TABLE fraud_rules (
    id SERIAL PRIMARY KEY,
    code CHARACTER VARYING(100),
    title CHARACTER VARYING(255),
    description CHARACTER VARYING(255),
    threshold NUMERIC(20, 0),
    enable BOOLEAN DEFAULT FALSE,
    severity CHARACTER VARYING(10),
    created_at TIMESTAMP default now()
);

comment on column fraud_rules.code is 'Техническое имя правила';
comment on column fraud_rules.title is 'Человекочитаемое название';
comment on column fraud_rules.description is 'Подробное описание логики';
comment on column fraud_rules.threshold is 'Пороговое значение (если применимо)';


INSERT INTO fraud_rules (code, title, description, threshold, enable, severity) VALUES
    ('high_amount', 'Лимит', 'Превышение лимита сумму', 500000, true, 'HIGH'),
    ('high_amount', 'Лимит', 'Превышение лимита сумму', 1000000, false, 'HIGH'),
    ('high_amount', 'Лимит', 'Превышение лимита сумму', 300000, false, 'HIGH'),

    ('geo_jump', 'VPN', 'Прыжок с 1 страны в другую', 60, true, 'LOW'),
    ('geo_jump', 'VPN', 'Превышение лимита сумму', 10, false, 'LOW'),
    ('geo_jump', 'VPN', 'Превышение лимита сумму', 600, false, 'LOW');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE fraud_rules;

-- +goose StatementEnd
