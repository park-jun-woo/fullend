CREATE TABLE reviews (
    id         BIGSERIAL    PRIMARY KEY,
    user_id    BIGINT       NOT NULL REFERENCES users(id),
    course_id  BIGINT       NOT NULL REFERENCES courses(id),
    rating     INT          NOT NULL,
    comment    TEXT         NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_review UNIQUE (user_id, course_id),
    CONSTRAINT chk_rating CHECK (rating >= 1 AND rating <= 5)
);

CREATE INDEX idx_reviews_course ON reviews (course_id, created_at);
CREATE INDEX idx_reviews_rating ON reviews (rating);
