package biz

import (
	"github.com/google/wire"
	"github.com/yeralin-munar/tt-go-json-fernet/internal/biz/usecase/datacollector"
	"github.com/yeralin-munar/tt-go-json-fernet/internal/biz/usecase/jsongenerator"
)

var ProviderSet = wire.NewSet(
	jsongenerator.NewJsonGeneratorUseCase,
	datacollector.NewDataCollectorUseCase,
)
