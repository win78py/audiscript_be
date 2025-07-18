package transcribe

type TranscribeRequest struct {
	Title     string `json:"title" binding:"required,min=3"`
	FileURL   string `json:"file_url" binding:"required,url"`
}