package schema

type CountFavoritesByWorkIDResponse struct {
	Total int `json:"total"`
}

type IsFavoriteResponse struct {
	IsFavorite bool `json:"is_favorite"`
}
