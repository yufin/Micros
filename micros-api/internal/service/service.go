package service

import (
	"github.com/google/wire"
	dwV2 "micros-api/internal/service/dw/v2"
	rcV2 "micros-api/internal/service/rc/v2"
)

// ProviderSet is service providers.
var ProviderSet = wire.
	NewSet(
		NewRcServiceServicer,
		NewRcRdmServiceServicer,
		NewTreeGraphServiceServicer,
		NewNetGraphServiceServicer,
		rcV2.NewRcServiceServicer,
		dwV2.NewDwServiceServicer,
	)
