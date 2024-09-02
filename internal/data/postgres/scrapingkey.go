package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/jackc/pgx/v5"
	"github.com/yeralin-munar/tt-go-json-fernet/internal/biz/dto"
	"github.com/yeralin-munar/tt-go-json-fernet/internal/data/dao"
)

type ScrapingKeyRepo struct {
	log *log.Helper

	db *DB
}

func NewScrapingKeyRepo(db *DB, logger log.Logger) *ScrapingKeyRepo {
	return &ScrapingKeyRepo{
		db:  db,
		log: log.NewHelper(logger),
	}
}

func (s *ScrapingKeyRepo) Query(
	ctx context.Context,
	filter *dto.ScrapingKeyFilter,
) (
	[]*dao.ScrapingKey,
	int32,
	error,
) {
	var (
		query, countQuery strings.Builder
	)
	query.WriteString(`
SELECT
	sk.id,
	sk.name,
	sk.created_at,
	sk.updated_at,
	sk.deleted_at
FROM scraping_key sk
`)

	var (
		conditions  []string
		queryParams = pgx.NamedArgs{}
	)

	if filter != nil && !filter.IncludesDeleted {
		conditions = append(conditions, "sk.deleted_at IS NULL ")
	}

	// For the count query
	var total int32

	if filter != nil && filter.HasTotalNumber {
		countQuery.WriteString(`
SELECT
	COUNT(1)
FROM scraping_key sk `)

		if len(conditions) > 0 {
			countQuery.WriteString(" WHERE ")
			countQuery.WriteString(strings.Join(conditions, " AND "))
		}

		err := s.db.GetConnection(ctx).
			QueryRow(ctx, countQuery.String(), queryParams).
			Scan(&total)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan count scraping keys: %w", err)
		}
	}

	// For the main query
	if filter != nil && filter.LastID > 0 {
		// TODO: Move Sub-query to WITH clause
		conditions = append(conditions, "sk.id > @LastID ")
		queryParams["LastID"] = filter.LastID
	}

	if len(conditions) > 0 {
		query.WriteString(" WHERE ")
		query.WriteString(strings.Join(conditions, " AND "))
	}

	query.WriteString(" ORDER BY sk.created_at ASC, sk.id ASC ")

	if filter != nil && filter.Limit > 0 {
		query.WriteString(fmt.Sprintf(" LIMIT %d ", filter.Limit))
	}

	if filter != nil && filter.HasLock {
		query.WriteString(" FOR UPDATE ")
	}

	rows, err := s.db.GetConnection(ctx).Query(ctx, query.String(), queryParams)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query scraping keys: %w", err)
	}

	defer rows.Close()

	res, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByPos[dao.ScrapingKey])
	if err != nil {
		return nil, 0, fmt.Errorf("failed to collect scraping keys rows: %w", err)
	}

	return res, total, nil
}
