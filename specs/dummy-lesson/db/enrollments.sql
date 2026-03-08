CREATE TABLE enrollments (
    id         BIGSERIAL    PRIMARY KEY,
    user_id    BIGINT       NOT NULL REFERENCES users(id),
    course_id  BIGINT       NOT NULL REFERENCES courses(id),
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_enrollment UNIQUE (user_id, course_id)
);

CREATE INDEX idx_enrollments_user   ON enrollments (user_id);
CREATE INDEX idx_enrollments_course ON enrollments (course_id);
