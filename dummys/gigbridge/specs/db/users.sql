CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL, -- @sensitive
    role VARCHAR(20) NOT NULL CHECK (role IN ('client', 'freelancer', 'admin')),
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO users (id, email, password_hash, role, name)
VALUES (0, 'nobody@system', '', 'system', 'Nobody')
ON CONFLICT DO NOTHING;
