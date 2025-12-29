-- create a user with set permissions 

CREATE OR REPLACE FUNCTION create_user(
    p_username TEXT,
    p_password TEXT,
    p_table_schemas TEXT[],
    p_sequence_schemas TEXT[],
    p_search_path TEXT DEFAULT NULL
)
RETURNS void
LANGUAGE plpgsql
AS $$
DECLARE schema_name TEXT;
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = p_username) THEN
        EXECUTE format('CREATE ROLES %I LOGIN PASSWORD %L', p_username, p_password);
    END IF;

    EXECUTE format('ALTER ROLE %I NOSUPERUSER NOCREATEDB NOCREATEROLE', p_username);

    FOREACH schema_name IN ARRAY p_table_schemas LOOP
        EXECUTE format('GRANT USAGE ON SCHEMA %I TO %I', schema_name, p_username);
    END LOOP;

    FOREACH schema_name IN ARRAY p_table_schemas LOOP
        EXECUTE format(
            'ALTER DEFAULT PRIVILEGES IN SCHEMA %I GRANT SELECT, INSERT, UPDATE ON TABLES TO %I',
            schema_name, p_username
        );
    END LOOP;

    FOREACH schema_name IN ARRAY p_sequence_schemas LOOP
        EXECUTE format(
            'ALTER DEFAULT PRIVILEGES IN SCHEMA %I GRANT USAGE, SELECT ON SEQUENCES TO %I',
            schema_name, p_username
        );
    END LOOP;

    F p_search_path IS NOT NULL THEN
        EXECUTE format('ALTER ROLE %I SET search_path = %s', p_username, p_search_path);
    END IF;
END;
$$