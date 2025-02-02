CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO users (email, password_hash, is_admin, created_at)
VALUES ('admin@example.com', '$2a$10$nLO909QZ5enK5xAfklkbMeXNAhxmt/GtJW2d3Ddx3U4D4IOYv0TEO', TRUE, CURRENT_TIMESTAMP);