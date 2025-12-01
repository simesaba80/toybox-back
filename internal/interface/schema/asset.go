package schema

type UploadAssetResponse struct {
	URL string `json:"url"`
}

func ToUploadAssetResponse(url string) UploadAssetResponse {
	return UploadAssetResponse{
		URL: url,
	}
}
