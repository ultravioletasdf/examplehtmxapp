PRAGMA foreign_keys = on;

CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    email TEXT NOT NULL,
    password BLOB NOT NULL,
    created_at INTEGER NOT NULL,
    verified INTEGER NOT NULL DEFAULT 0,
    verify_code TEXT,
    verify_expire_at INTEGER
);

CREATE UNIQUE INDEX idx_users_email ON users (email);

CREATE TABLE sessions (
    token TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    expire_at INTEGER NOT NULL,
    created_at INTEGER NOT NULL
);
