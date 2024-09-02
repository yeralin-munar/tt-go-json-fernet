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

type FileDataRepo struct {
	log *log.Helper

	db *DB
}

func NewFileDataRepo(db *DB, logger log.Logger) *FileDataRepo {
	return &FileDataRepo{
		db:  db,
		log: log.NewHelper(logger),
	}
}

func (f *FileDataRepo) Insert(
	ctx context.Context,
	fileDatas []*dao.FileData,
) (
	[]*dao.FileData,
	error,
) {
	query := `
INSERT INTO file_data (
	name,
	data,
	created_at,
	updated_at
)
SELECT
	fd.name,
	fd.data,
	NOW(),
	NOW()
FROM UNNEST(@FileDatas::file_data_type[]) fd
RETURNING
	id,
	name,
	data,
	created_at,
	updated_at,
	deleted_at
`
	rows, err := f.db.GetConnection(ctx).
		Query(
			ctx,
			query,
			pgx.NamedArgs{
				"FileDatas": fileDatas,
			},
		)
	if err != nil {
		return nil, fmt.Errorf("failed to insert file datas: %w", err)
	}

	defer rows.Close()

	res, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByNameLax[dao.FileData])
	if err != nil {
		return nil, fmt.Errorf("failed to collect file data rows: %w", err)
	}

	return res, nil
}

func (f *FileDataRepo) Query(
	ctx context.Context,
	filter *dto.FileDataFilter,
) (
	[]*dao.FileData,
	int32,
	error,
) {
	var (
		query, countQuery strings.Builder
	)
	query.WriteString(`
SELECT
	fd.id,
	fd.name,
	fd.data,
	fd.created_at,
	fd.updated_at,
	fd.deleted_at
FROM file_data fd
`)

	var (
		conditions  []string
		queryParams = pgx.NamedArgs{}
	)

	if filter != nil && !filter.IncludesDeleted {
		conditions = append(conditions, "deleted_at IS NULL ")
	}

	// For the count query
	var total int32

	if filter != nil && filter.HasTotalNumber {
		countQuery.WriteString(`
SELECT
	COUNT(1)
FROM file_data fd `)

		if len(conditions) > 0 {
			countQuery.WriteString(" WHERE ")
			countQuery.WriteString(strings.Join(conditions, " AND "))
		}

		err := f.db.GetConnection(ctx).
			QueryRow(ctx, countQuery.String(), queryParams).
			Scan(&total)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan count file datas: %w", err)
		}
	}

	// For the main query
	if filter != nil && filter.LastID > 0 {
		// TODO: Move Sub-query to WITH clause
		conditions = append(conditions, fmt.Sprintf("fd.id > %d ", filter.LastID))
		queryParams["LastID"] = filter.LastID
	}

	if len(conditions) > 0 {
		query.WriteString(" WHERE ")
		query.WriteString(strings.Join(conditions, " AND "))
	}

	query.WriteString(" ORDER BY fd.created_at ASC, fd.id ASC ")

	if filter != nil && filter.Limit > 0 {
		query.WriteString(fmt.Sprintf(" LIMIT %d ", filter.Limit))
	}

	if filter != nil && filter.HasLock {
		query.WriteString(" FOR UPDATE ")
	}

	rows, err := f.db.GetConnection(ctx).Query(ctx, query.String(), queryParams)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query file datas: %w", err)
	}

	defer rows.Close()

	res, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByPos[dao.FileData])
	if err != nil {
		return nil, 0, fmt.Errorf("failed to collect file data rows: %w", err)
	}

	return res, total, nil
}
