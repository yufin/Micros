package dto

type Node struct {
	Vid        int64                  `json:"vid"`
	Labels     []string               `json:"labels"`
	Properties map[string]interface{} `json:"properties"`
}
