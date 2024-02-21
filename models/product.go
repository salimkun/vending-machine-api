package models

import (
	"time"

	"github.com/guregu/null"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name      string      `json:"name" gorm:"unique"`
	Price     int         `json:"price"`
	CreatedBy string      `json:"created_by"`
	UpdatedBy string      `json:"updated_by"`
	DeletedBy null.String `json:"deleted_by"`
}

type CreateProductInput struct {
	Name  string `json:"name" binding:"required"`
	Price int    `json:"price" binding:"required"`
}

type UpdateProductInput struct {
	Name      string    `json:"name"`
	Price     int       `json:"price"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
}

type DeleteProductInput struct {
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	DeletedAt time.Time `json:"deleted_at"`
	DeletedBy string    `json:"deleted_by"`
}
