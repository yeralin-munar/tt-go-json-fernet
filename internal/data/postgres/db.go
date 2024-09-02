package postgres

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yeralin-munar/tt-go-json-fernet/config"
)

type ctxDBKey string

const dbKey ctxDBKey = "postgres"

var customTypeMap = make(map[string]*pgtype.Type)

type DB struct {
	pool *pgxpool.Pool
}

func GenerateDBURL(data *config.Data) string {
	return fmt.Sprintf(
		// "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"postgres://%s:%s@%s:%d/%s",
		data.Database.User,
		data.Database.Password,
		data.Database.Host,
		data.Database.Port,
		data.Database.Name,
	)
}

func NewDB(data *config.Data) (*DB, error) {
	cfg, err := pgxpool.ParseConfig(GenerateDBURL(data))
	if err != nil {
		return nil, err
	}

	// Register for custom data types
	cfg.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		return registerDataTypes(ctx, conn)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, err
	}

	// Ping to check if the connection is successful
	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	return &DB{
		pool: pool,
	}, nil
}

func (d *DB) GetConnection(ctx context.Context) *pgxpool.Pool {
	if v, ok := ctx.Value(dbKey).(*pgxpool.Pool); ok {
		return v
	}

	return d.pool
}

func registerDataTypes(ctx context.Context, conn *pgx.Conn) error {
	customTypeNames := []string{
		"scraping_key_type",
		"scraping_key_type[]",
		"file_data_type",
		"file_data_type[]",
	}

	for _, typeName := range customTypeNames {
		dataType, ok := customTypeMap[typeName]
		if !ok {
			// Load the type from the database
			loadedType, err := conn.LoadType(ctx, typeName)
			if err != nil {
				log.Errorf("Failed to load type %s: %v", typeName, err)

				return err
			}
			dataType = loadedType

			// Cache the type for future use
			customTypeMap[typeName] = dataType
		}

		conn.TypeMap().RegisterType(dataType)
	}

	return nil
}
