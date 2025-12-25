-- +goose Up
-- +goose StatementBegin
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    transaction_id UUID DEFAULT gen_random_uuid(),
    account_id NUMERIC(20, 0) CHECK (account_id >= 0),
    amount NUMERIC(20, 0) CHECK (amount >= 0),
    currency_id INTEGER REFERENCES currencies (id),
    merchant CHARACTER VARYING(255),
    country_id INTEGER REFERENCES countries (id),
    status_id INTEGER REFERENCES statuses (id),
    payload JSONB,
    source_id INTEGER REFERENCES sources (id),
    created_at TIMESTAMP default now(),
    ingested_at TIMESTAMP
);

CREATE INDEX idx_transactions_account_id ON transactions(account_id);
CREATE INDEX idx_transactions_created_at ON transactions(created_at);
CREATE INDEX idx_transactions_currency_id ON transactions(currency_id);
CREATE INDEX idx_transactions_country_id ON transactions(country_id);
CREATE INDEX idx_transactions_status_id ON transactions(status_id);
CREATE INDEX idx_transactions_source_id ON transactions(source_id);
CREATE INDEX idx_transactions_ingested_at ON transactions(ingested_at);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE transactions;
-- +goose StatementEnd
