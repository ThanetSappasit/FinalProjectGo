package dto

type ProductSearchRequest struct {
	ProductName *string `json:"pname"`
	MinPrice    *string `json:"min_price"`
	MaxPrice    *string `json:"max_price"`
	Description *string `json:"description"`
}
