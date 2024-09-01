//go:build wireinject
// +build wireinject

package main

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/yeralin-munar/tt-go-json-fernet/cmd"
	"github.com/yeralin-munar/tt-go-json-fernet/config"
	"github.com/yeralin-munar/tt-go-json-fernet/internal/biz"
	"github.com/yeralin-munar/tt-go-json-fernet/internal/biz/usecase"
	"github.com/yeralin-munar/tt-go-json-fernet/internal/biz/usecase/datacollector"
	"github.com/yeralin-munar/tt-go-json-fernet/internal/biz/usecase/jsongenerator"
	"github.com/yeralin-munar/tt-go-json-fernet/internal/data"
	"github.com/yeralin-munar/tt-go-json-fernet/internal/data/postgres"
	"github.com/yeralin-munar/tt-go-json-fernet/internal/encryption"
	"github.com/yeralin-munar/tt-go-json-fernet/internal/encryption/fernet"
	"github.com/yeralin-munar/tt-go-json-fernet/internal/service"
)

func wireApp(log.Logger, *config.Server, *config.Data) (*cmd.App, error) {
	panic(
		wire.Build(
			service.ProviderSet,
			biz.ProviderSet,
			data.ProviderSet,
			encryption.ProviderSet,
			wire.Bind(new(usecase.TransactionManager), new(*postgres.TransactionManager)),
			wire.Bind(new(usecase.ScrapingKeyRepo), new(*postgres.ScrapingKeyRepo)),
			wire.Bind(new(usecase.FileDataRepo), new(*postgres.FileDataRepo)),
			wire.Bind(new(usecase.Cipherator), new(*fernet.FernetCipherator)),
			wire.Bind(new(service.JsonGeneratorUseCase), new(*jsongenerator.JsonGeneratorUseCase)),
			wire.Bind(new(service.DataCollectorUseCase), new(*datacollector.DataCollectorUseCase)),
			cmd.NewApp,
		))
}
