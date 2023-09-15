package service

import (
	"github.com/google/wire"
	dwV2 "micros-dw/internal/service/dw/v2"
)

// ProviderSet is service providers.
var ProviderSet = wire.
	NewSet(
		dwV2.NewDwServiceServicer,
	)
