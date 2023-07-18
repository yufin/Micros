package biz

import "github.com/google/wire"

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	NewOssMetadataUsecase,
	NewRcReportOssUsecase,
	NewGraphNodeUsecase,
	NewRcProcessedContentUsecase,
	NewRcOriginContentUsecase,
	NewRcDependencyDataUsecase,
)
