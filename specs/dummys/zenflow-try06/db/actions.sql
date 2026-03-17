CREATE TABLE actions (
    id BIGSERIAL PRIMARY KEY,
    workflow_id BIGINT NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    action_type VARCHAR(100) NOT NULL,
    payload_template TEXT NOT NULL DEFAULT '',
    sequence_order BIGINT NOT NULL
);

CREATE INDEX idx_actions_workflow_id ON actions(workflow_id);
