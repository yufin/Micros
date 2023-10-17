package dto

type PaginationReq struct {
	PageNum  int `json:"page_num"`
	PageSize int `json:"page_size"`
}

type PaginationResp struct {
	Total     int64       `json:"total"`
	TotalPage int         `json:"total_page"`
	PageNum   int         `json:"page_num"`
	PageSize  int         `json:"page_size"`
	Data      interface{} `json:"data"`
}
