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
DECLARE
    schema_name TEXT;
    obj_name TEXT;
BEGIN
    -- Create role if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = p_username) THEN
        EXECUTE format('CREATE ROLE %I LOGIN PASSWORD %L', p_username, p_password);
    END IF;

    EXECUTE format('ALTER ROLE %I NOSUPERUSER NOCREATEDB NOCREATEROLE', p_username);

    -- Grant schema usage
    FOREACH schema_name IN ARRAY p_table_schemas LOOP
        EXECUTE format('GRANT USAGE ON SCHEMA %I TO %I', schema_name, p_username);

        -- Grant privileges on existing tables in schema
        FOR obj_name IN
            SELECT tablename FROM pg_tables WHERE schemaname = schema_name
        LOOP
            EXECUTE format('GRANT SELECT, INSERT, UPDATE ON TABLE %I.%I TO %I',
                           schema_name, obj_name, p_username);
        END LOOP;

        -- Set default privileges for future tables
        EXECUTE format(
            'ALTER DEFAULT PRIVILEGES IN SCHEMA %I GRANT SELECT, INSERT, UPDATE ON TABLES TO %I',
            schema_name, p_username
        );
    END LOOP;

    -- Grant privileges on existing sequences and set default privileges
    FOREACH schema_name IN ARRAY p_sequence_schemas LOOP
        FOR obj_name IN
            SELECT sequencename FROM pg_sequences WHERE schemaname = schema_name
        LOOP
            EXECUTE format('GRANT USAGE, SELECT ON SEQUENCE %I.%I TO %I',
                           schema_name, obj_name, p_username);
        END LOOP;

        EXECUTE format(
            'ALTER DEFAULT PRIVILEGES IN SCHEMA %I GRANT USAGE, SELECT ON SEQUENCES TO %I',
            schema_name, p_username
        );
    END LOOP;

    -- Set search_path if provided
    IF p_search_path IS NOT NULL THEN
        EXECUTE format('ALTER ROLE %I SET search_path = %L', p_username, p_search_path);
    END IF;
END;
$$;