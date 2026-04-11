CREATE TABLE transactions (
    id BIGSERIAL PRIMARY KEY,
    gig_id BIGINT NOT NULL REFERENCES gigs(id),
    tx_type VARCHAR(20) NOT NULL CHECK (tx_type IN ('hold', 'release', 'refund')),
    amount INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_transactions_gig_id ON transactions (gig_id);
