package encryption

import (
	"github.com/google/wire"
	"github.com/yeralin-munar/tt-go-json-fernet/internal/encryption/fernet"
)

var ProviderSet = wire.NewSet(
	fernet.NewFernetCipherator,
)
