CREATE TABLE lessons (
    id         BIGSERIAL    PRIMARY KEY,
    course_id  BIGINT       NOT NULL REFERENCES courses(id),
    title      VARCHAR(200) NOT NULL,
    video_url  VARCHAR(500) NOT NULL DEFAULT '',
    sort_order INT          NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW() -- @archived
);

CREATE INDEX idx_lessons_course ON lessons (course_id, sort_order);
