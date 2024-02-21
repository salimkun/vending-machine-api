package payload

type BuyProductRequest struct {
	Money []int `json:"money" validate:"min=1,dive,oneof=2000 5000"`
}
