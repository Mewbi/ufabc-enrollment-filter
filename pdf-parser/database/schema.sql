-- Table to Enrollments
CREATE TABLE IF NOT EXISTS enrollments (
    id          TEXT PRIMARY KEY,
    name        TEXT,
    url         TEXT,
    content     BLOB,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);
