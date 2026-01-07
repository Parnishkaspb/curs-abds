CREATE TABLE transactions (
                                     transaction_id String,
                                     created_at     DateTime64(3),
                                     account_id     UInt64,
                                     amount         UInt64,
                                     country        String,
                                     merchant       String,
                                     accepted       Bool
)
    ENGINE MergeTree()
ORDER BY (created_at);