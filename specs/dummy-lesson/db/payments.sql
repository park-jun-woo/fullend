CREATE TABLE payments (
    id            BIGSERIAL    PRIMARY KEY,
    user_id       BIGINT       NOT NULL REFERENCES users(id),
    enrollment_id BIGINT       NOT NULL REFERENCES enrollments(id),
    amount        INT          NOT NULL,
    payment_method VARCHAR(20) NOT NULL,
    status        VARCHAR(20)  NOT NULL DEFAULT 'pending',
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payments_user       ON payments (user_id, created_at);
CREATE INDEX idx_payments_enrollment ON payments (enrollment_id);
