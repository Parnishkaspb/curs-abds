-- +goose Up
-- +goose StatementBegin
CREATE TABLE transcations (
    id SERIAL,
    transaction_id UUID DEFAULT gen_random_uuid(),
    account_id NUMERIC(20, 0) CHECK (id >= 0),
    amount NUMERIC(20, 0) CHECK (id >= 0),
    currency_id INTEGER REFERENCES currencies (id),
    merchant CHARACTER VARYING(255),
    country_id INTEGER REFERENCES countries (id),
    status_id INTEGER REFERENCES statuses (id),
    payload JSONB,
    source_id INTEGER REFERENCES sources (id),
    created_at TIMESTAMP default now(),
    ingested_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE transcations;
-- +goose StatementEnd
