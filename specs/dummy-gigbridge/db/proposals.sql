CREATE TABLE proposals (
    id BIGSERIAL PRIMARY KEY,
    gig_id BIGINT NOT NULL REFERENCES gigs(id),
    freelancer_id BIGINT NOT NULL REFERENCES users(id),
    bid_amount INT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_proposals_gig_id ON proposals (gig_id);
CREATE INDEX idx_proposals_freelancer_id ON proposals (freelancer_id);
CREATE INDEX idx_proposals_status ON proposals (status);
