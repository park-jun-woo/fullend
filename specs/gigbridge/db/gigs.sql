CREATE TABLE gigs (
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NOT NULL REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    budget INTEGER NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'draft',
    freelancer_id BIGINT NOT NULL DEFAULT 0 REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_gigs_client_id ON gigs (client_id);
CREATE INDEX idx_gigs_freelancer_id ON gigs (freelancer_id);
CREATE INDEX idx_gigs_status ON gigs (status);
CREATE INDEX idx_gigs_budget ON gigs (budget);
CREATE INDEX idx_gigs_created_at ON gigs (created_at);
