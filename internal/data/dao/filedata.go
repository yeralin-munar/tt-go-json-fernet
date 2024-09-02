package dao

import "time"

type FileData struct {
	ID        int64                  `db:"id"`
	Name      string                 `db:"name"`
	Data      map[string]interface{} `db:"data"`
	CreatedAt time.Time              `db:"created_at"`
	UpdatedAt time.Time              `db:"updated_at"`
	DeletedAt *time.Time             `db:"deleted_at"`
}
