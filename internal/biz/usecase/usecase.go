package usecase

import (
	"context"

	"github.com/yeralin-munar/tt-go-json-fernet/internal/biz/dto"
	"github.com/yeralin-munar/tt-go-json-fernet/internal/data/dao"
)

type TransactionManager interface {
	InTransaction(ctx context.Context, txFunc func(ctx context.Context) error) error
}

type ScrapingKeyRepo interface {
	Query(ctx context.Context, filter *dto.ScrapingKeyFilter) ([]*dao.ScrapingKey, int32, error)
}

type FileDataRepo interface {
	Insert(ctx context.Context, fileDatas []*dao.FileData) ([]*dao.FileData, error)
	Query(ctx context.Context, filter *dto.FileDataFilter) ([]*dao.FileData, int32, error)
}

type Cipherator interface {
	Encrypt(msg []byte, val string) ([]byte, error)
	Decrypt(msg []byte, val string) ([]byte, error)
}
