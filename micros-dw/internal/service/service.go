package service

import (
	"github.com/google/wire"
	v2 "micros-dw/internal/service/v2"
	v3 "micros-dw/internal/service/v3"
)

// ProviderSet is service providers.
var ProviderSet = wire.
	NewSet(
		//NewDwdataServiceServicer,
		v2.NewDwdataServiceServicer,
		v3.NewDwdataServiceServicer,
	)
