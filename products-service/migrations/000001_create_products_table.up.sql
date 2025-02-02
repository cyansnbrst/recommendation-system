CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    tags TEXT[] NOT NULL,
    version integer NOT NULL DEFAULT 1
);
