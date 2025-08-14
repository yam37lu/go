package response

type PageData struct {
	Total int64       `json:"total"`
	Data  interface{} `json:"data,omitempty"`
}

type ResData struct {
	Data interface{} `json:"data,omitempty"`
}
