DO
$$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'user') THEN
        CREATE USER "user" WITH PASSWORD 'password';
        ALTER USER "user" CREATEDB;
    END IF;
END
$$;

CREATE DATABASE users OWNER "user";
GRANT ALL PRIVILEGES ON DATABASE users TO "user";

CREATE DATABASE profiles OWNER "user";
GRANT ALL PRIVILEGES ON DATABASE profiles TO "user";

CREATE DATABASE products OWNER "user";
GRANT ALL PRIVILEGES ON DATABASE products TO "user";

CREATE DATABASE recommendations OWNER "user";
GRANT ALL PRIVILEGES ON DATABASE recommendations TO "user";

CREATE DATABASE analytics OWNER "user";
GRANT ALL PRIVILEGES ON DATABASE analytics TO "user";
