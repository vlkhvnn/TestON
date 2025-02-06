DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'teston') THEN
        CREATE DATABASE teston;
    END IF;
END
$$;
