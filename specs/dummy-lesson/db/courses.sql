CREATE TABLE courses (
    id            BIGSERIAL    PRIMARY KEY,
    instructor_id BIGINT       NOT NULL REFERENCES users(id),
    title         VARCHAR(200) NOT NULL,
    description   TEXT         NOT NULL DEFAULT '',
    category      VARCHAR(50)  NOT NULL,
    level         VARCHAR(20)  NOT NULL DEFAULT 'beginner',
    price         INT          NOT NULL DEFAULT 0,
    published     BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_courses_instructor ON courses (instructor_id);
CREATE INDEX idx_courses_category   ON courses (category);
CREATE INDEX idx_courses_created    ON courses (created_at);
CREATE INDEX idx_courses_price      ON courses (price);
