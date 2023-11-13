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

type DeleteFilesResponse struct {
	Count  int
	Result []DeleteFilesResponseItem
}

type DeleteFilesResponseItem struct {
	ID    string `json:"id"`
	Error string `json:"error,omitempty"`
}

type SearchFilesResponse struct {
	Count int               `json:"count"`
	Files []GetFileResponse `json:"files"`
}
