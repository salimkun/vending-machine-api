package payload

import "time"

type CreateProductRequest struct {
	Name  string `json:"name" validate:"required,min=1,max=20"`
	Price int    `json:"price" validate:"gte=2000,lte=50000"`
}

type UpdateProductRequest struct {
	Name  string `json:"name" validate:"omitempty,min=1,max=20"`
	Price int    `json:"price" validate:"omitempty,gte=2000,lte=50000"`
}

type ProductResponse struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name" gorm:"unique"`
	Price     int       `json:"price"`
	CreatedBy string    `json:"created_by"`
	UpdatedBy string    `json:"updated_by"`
}
