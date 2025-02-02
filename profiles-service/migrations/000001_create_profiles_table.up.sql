CREATE TABLE profiles (
    user_uid TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    location TEXT NOT NULL,
    interests TEXT[] NOT NULL
);
