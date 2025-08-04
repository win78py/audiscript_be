package transcribe

type CreateAudioRequest struct {
	Title     string `json:"title" binding:"required,min=3"`
	FileURL   string `json:"file_url" binding:"required,url"`
}

type TranscribeRequest struct {
    FileURL  string `json:"file_url" binding:"required"`
    Language string `json:"language" binding:"required"`
    AudioID  string `json:"audio_id" binding:"required"`
}