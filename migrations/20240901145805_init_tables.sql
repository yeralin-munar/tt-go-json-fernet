-- +goose Up
-- +goose StatementBegin
-- Table: file_data
CREATE TABLE file_data (
  id SERIAL PRIMARY KEY,
  name TEXT,
  data JSONB NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMP
);

DO
$$ 
  BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'file_data_type') THEN
      CREATE TYPE file_data_type AS
      (
        id BIGINT,
        name TEXT,
        data JSONB,
        created_at  TIMESTAMP WITH TIME ZONE,
        updated_at  TIMESTAMP WITH TIME ZONE,
        deleted_at  TIMESTAMP WITH TIME ZONE
      );
    END IF;
  END;
$$
LANGUAGE plpgsql;

-- Table scraping_key
CREATE TABLE scraping_key (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMP
);

DO
$$ 
  BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'scraping_key_type') THEN
      CREATE TYPE scraping_key_type AS
      (
        id BIGINT,
        name TEXT,
        created_at  TIMESTAMP WITH TIME ZONE,
        updated_at  TIMESTAMP WITH TIME ZONE,
        deleted_at  TIMESTAMP WITH TIME ZONE
      );
    END IF;
  END;
$$
LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE file_data;
DROP TYPE IF EXISTS setting_type;
DROP TABLE scraping_key;
DROP TYPE IF EXISTS scraping_key_type;
-- +goose StatementEnd
