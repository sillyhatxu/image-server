package dto

type MultiFileDTO struct {
	UploadFile string `json:"upload_file"`

	OutputFile string `json:"output_file"`
}

type UploadURLDTO struct {
	URL string `json:"url"`
}
