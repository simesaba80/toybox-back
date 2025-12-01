package schema

import "github.com/simesaba80/toybox-back/internal/domain/entity"

type UploadAssetResponse struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

func ToUploadAssetResponse(asset *entity.Asset) UploadAssetResponse {
	return UploadAssetResponse{
		ID:  asset.ID,
		URL: asset.URL,
	}
}
