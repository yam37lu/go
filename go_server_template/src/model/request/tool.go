package request

import "mime/multipart"

type ExportToolRequest struct {
	Table    string `json:"table"`
	Where    string `json:"where"`
	Truncate bool   `json:"truncate"`
	FileName string `json:"fileName"`
	FilePath string `json:"filePath"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	Download bool   `json:"download"`
}

type ImportToolRequest struct {
	FilePath string                `json:"filePath"`
	FileName string                `json:"fileName"`
	File     *multipart.FileHeader `json:"-"`
}

type ImportCount struct {
	Table string `json:"table"`
	Where string `json:"where"`
}
