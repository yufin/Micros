package biz

import "github.com/google/wire"

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	NewOssMetadataUsecase,
	NewRcReportOssUsecase,
	NewGraphUsecase,
	NewRcProcessedContentUsecase,
	NewRcOriginContentUsecase,
	NewRcDependencyDataUsecase,
	NewRcRdmResultUsecase,
	NewRcRdmResDetailUsecase,
	NewMgoRcUsecase,
	NewDwEnterpriseUsecase,
	NewRcDecisionFactorUsecase,
)
