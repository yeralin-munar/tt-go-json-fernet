package dto

import (
	"time"

	"github.com/yeralin-munar/tt-go-json-fernet/internal/data/dao"
)

type ScrapingKey struct {
	ID        int64      `json:"id,omitempty"`
	Name      string     `json:"name,omitempty"`
	CreatedAt time.Time  `json:"created_at,omitempty"`
	UpdatedAt time.Time  `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func ConvertScrapingKeyDaoToDto(dao *dao.ScrapingKey) *ScrapingKey {
	return &ScrapingKey{
		ID:        dao.ID,
		Name:      dao.Name,
		CreatedAt: dao.CreatedAt,
		UpdatedAt: dao.UpdatedAt,
		DeletedAt: dao.DeletedAt,
	}
}

func ConvertScrapingKeyDaosToDtos(daos []*dao.ScrapingKey) []*ScrapingKey {
	dtos := make([]*ScrapingKey, len(daos))
	for i, scDAO := range daos {
		dtos[i] = ConvertScrapingKeyDaoToDto(scDAO)
	}
	return dtos
}

type ScrapingKeyFilter struct {
	LastID          int64
	Limit           int32
	IncludesDeleted bool
	HasTotalNumber  bool
	HasLock         bool
}
