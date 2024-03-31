package service

import (
	"github.com/google/wire"
	dwV2 "micros-api/internal/service/dw/v2"
	dwV3 "micros-api/internal/service/dw/v3"
	rcV3 "micros-api/internal/service/rc/v3"
)

// ProviderSet is service providers.
var ProviderSet = wire.
	NewSet(
		NewRcServiceServicer,
		NewRcRdmServiceServicer,
		NewTreeGraphServiceServicer,
		NewNetGraphServiceServicer,
		dwV2.NewDwServiceServicer,
		dwV3.NewDwServiceServicer,
		rcV3.NewRcServiceServicer,
	)
