CREATE TABLE actions (
    id BIGSERIAL PRIMARY KEY,
    workflow_id BIGINT NOT NULL REFERENCES workflows(id),
    action_type VARCHAR(50) NOT NULL,
    payload_template JSONB NOT NULL DEFAULT '{}',
    sequence_order INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_actions_workflow_id ON actions (workflow_id);
