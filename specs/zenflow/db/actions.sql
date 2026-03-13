CREATE TABLE actions (
    id BIGSERIAL PRIMARY KEY,
    workflow_id BIGINT NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    action_type TEXT NOT NULL,
    payload_template JSONB NOT NULL DEFAULT '{}',
    sequence_order INTEGER NOT NULL
);
