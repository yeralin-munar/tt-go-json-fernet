package dto

import (
	"time"

	"github.com/yeralin-munar/tt-go-json-fernet/internal/data/dao"
)

type FileData struct {
	ID        int64                  `json:"id,omitempty"`
	Name      string                 `json:"name,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
	CreatedAt time.Time              `json:"created_at,omitempty"`
	UpdatedAt time.Time              `json:"updated_at,omitempty"`
	DeletedAt *time.Time             `json:"deleted_at,omitempty"`
}

type FileDataFilter struct {
	LastID          int64
	Limit           int32
	IncludesDeleted bool
	HasTotalNumber  bool
	HasLock         bool
}

func ConvertFileDataDtosToDao(dtos []*FileData) []*dao.FileData {
	daos := make([]*dao.FileData, len(dtos))
	for i, fdDTO := range dtos {
		daos[i] = ConvertFileDataDtoToDao(fdDTO)
	}
	return daos
}

func ConvertFileDataDtoToDao(dto *FileData) *dao.FileData {
	return &dao.FileData{
		ID:        dto.ID,
		Name:      dto.Name,
		Data:      dto.Data,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
		DeletedAt: dto.DeletedAt,
	}
}
