package data

import (
	"github.com/google/wire"
	"github.com/yeralin-munar/tt-go-json-fernet/internal/data/postgres"
)

var ProviderSet = wire.NewSet(
	postgres.NewDB,
	postgres.NewTransactionManager,
	postgres.NewFileDataRepo,
	postgres.NewScrapingKeyRepo,
)
