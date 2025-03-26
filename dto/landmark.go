package dto

type Landmark struct {
	Name    string `json:"name"`
	Country int    `json:"country"`
	Detail  string `json:"detail"`
	Url     string `json:"url"`
}

type UpdateLandmarkRequest struct {
	Idx     int     `json:"idx" binding:"required"`
	Name    *string `json:"name"`
	Country *int    `json:"country"`
	Detail  *string `json:"detail"`
	Url     *string `json:"url"`
}
