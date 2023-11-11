package dto

type UploadFileResponse struct {
	ID string `json:"id"`
}

type GetFileResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Size int64  `json:"size"`
}

type DeleteFileResponse struct {
	ID string `json:"id"`
}
