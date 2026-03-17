-- @archived
CREATE TABLE transactions (
    id BIGSERIAL PRIMARY KEY,
    gig_id BIGINT NOT NULL REFERENCES gigs(id),
    tx_type VARCHAR(50) NOT NULL,
    amount INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_transactions_gig_id ON transactions(gig_id);
