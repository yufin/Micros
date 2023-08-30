package service

import (
	"brillinkmicros/internal/service/rc/v2"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.
	NewSet(
		NewRcServiceServicer,
		NewRcRdmServiceServicer,
		NewTreeGraphServiceServicer,
		NewNetGraphServiceServicer,
		v2.NewRcServiceServicer,
	)
