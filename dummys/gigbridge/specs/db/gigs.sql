CREATE TABLE gigs (
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NOT NULL REFERENCES users(id),
    title VARCHAR(200) NOT NULL,
    description TEXT NOT NULL,
    budget INTEGER NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'open', 'in_progress', 'under_review', 'completed', 'disputed')),
    freelancer_id BIGINT NOT NULL DEFAULT 0 REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_gigs_client_id ON gigs (client_id);
CREATE INDEX idx_gigs_status ON gigs (status);
CREATE INDEX idx_gigs_created_at ON gigs (created_at);
