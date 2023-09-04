package service

import (
	dwV2 "brillinkmicros/internal/service/dw/v2"
	rcV2 "brillinkmicros/internal/service/rc/v2"
	"github.com/google/wire"
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
