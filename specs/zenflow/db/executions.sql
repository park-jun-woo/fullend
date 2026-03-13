CREATE TABLE executions (
    id BIGSERIAL PRIMARY KEY,
    workflow_id BIGINT NOT NULL REFERENCES workflows(id),
    org_id BIGINT NOT NULL REFERENCES organizations(id),
    status TEXT NOT NULL DEFAULT 'pending',
    credits_spent INTEGER NOT NULL DEFAULT 0,
    executed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
