-- +goose Up
-- +goose StatementBegin
CREATE TABLE fraud_alerts (
    id SERIAL PRIMARY KEY,
    transaction_id INTEGER REFERENCES transactions (id),
    account_id INTEGER,
    fraud_rule INTEGER REFERENCES fraud_rules (id),
    description CHARACTER VARYING(255),
    created_at TIMESTAMP default now(),
    resolved BOOLEAN DEFAULT FALSE
);


-- CREATE INDEX idx_fraud_alerts_account_id ON fraud_alerts(account_id);
-- CREATE INDEX idx_fraud_alerts_created_at ON fraud_alerts(created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE fraud_alerts;
-- +goose StatementEnd

