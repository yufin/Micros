package dto

type Edge struct {
	SourceId   int64                  `json:"source_id"`
	TargetId   int64                  `json:"to"`
	Label      string                 `json:"label"`
	Rank       int64                  `json:"rank"`
	Properties map[string]interface{} `json:"properties"`
}
