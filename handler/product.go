package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
	"github.com/salimkun/vending-machine-api/database"
	"github.com/salimkun/vending-machine-api/middleware"
	"github.com/salimkun/vending-machine-api/models"
	"github.com/salimkun/vending-machine-api/payload"
)

// @Summary Get Product
// @Description Product List
// @Accept  json
// @Produce  json
// @Tags Product
// @Router /api/product [get]
// @Success 200 {object} payload.Response200{data=[]payload.ProductResponse}
// @Success 400 {object} payload.ResponseError
// @Success 404 {object} payload.ResponseError
// @Success 401 {object} payload.ResponseError
// @Success 500 {object} payload.ResponseError
// @Security  ApiJwtToken
func GetAllProduct(c *gin.Context) {
	var products []models.Product
	var responses []payload.ProductResponse
	database.Instance.Find(&products)
	_ = copier.Copy(&responses, &products)

	c.JSON(http.StatusOK, payload.Response200{
		Success: true,
		Message: "success get all data",
		Data:    responses,
	})
}

// @Summary Get Detai Product
// @Description Product detail
// @Accept  json
// @Produce  json
// @Tags Product
// @Param id path string true "ID Product"
// @Router /api/product/{id} [get]
// @Success 200 {object} payload.Response200{data=payload.ProductResponse}
// @Success 400 {object} payload.ResponseError
// @Success 404 {object} payload.ResponseError
// @Success 401 {object} payload.ResponseError
// @Success 500 {object} payload.ResponseError
// @Security  ApiJwtToken
func GetDetailProductByID(c *gin.Context) {
	var product models.Product
	var response payload.ProductResponse

	if err := database.Instance.Where("id = ?", c.Param("id")).First(&product).Error; err != nil {
		c.JSON(http.StatusBadRequest,
			payload.ResponseError{
				Success: true,
				Message: "error",
				Error:   "record not found",
			})
		return
	}

	_ = copier.Copy(&response, &product)

	c.JSON(http.StatusOK, payload.Response200{
		Success: true,
		Message: "success get detail data",
		Data:    product,
	})
}

// @Summary Create Product
// @Description Create Product
// @Accept  json
// @Produce  json
// @Tags Product
// @Router /api/product [post]
// @Param payload body payload.CreateProductRequest true "create product payload"
// @Success 200 {object} payload.Response200
// @Success 400 {object} payload.ResponseError
// @Success 404 {object} payload.ResponseError
// @Success 401 {object} payload.ResponseError
// @Success 500 {object} payload.ResponseError
// @Security  ApiJwtToken
func CreateProduct(c *gin.Context) {
	var request payload.CreateProductRequest

	tokenString := c.GetHeader("Authorization")
	tokenString = strings.ReplaceAll(tokenString, "Bearer ", "")
	claims, err := middleware.GetUserNameByToken(tokenString)
	if err != nil {
		c.JSON(http.StatusBadRequest, payload.ResponseError{
			Success: true,
			Message: "error",
			Error:   err.Error(),
		})
		return
	}

	err = json.NewDecoder(c.Request.Body).Decode(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, payload.ResponseError{
			Success: true,
			Message: "error",
			Error:   "invalid JSON",
		})
		return
	}

	validate := validator.New()
	// Validate struct
	err = validate.Struct(request)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		c.JSON(http.StatusBadRequest, payload.ResponseError{
			Success: true,
			Message: "error",
			Error:   fmt.Sprintf("%s", errors),
		})
		return
	}

	// Create product
	product := models.Product{Name: request.Name, Price: request.Price, CreatedBy: claims.Email, UpdatedBy: claims.Email}
	record := database.Instance.Create(&product)
	if record.Error != nil {
		c.JSON(http.StatusInternalServerError,
			payload.ResponseError{
				Success: true,
				Message: "error",
				Error:   err.Error(),
			})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, payload.Response200{
		Success: true,
		Message: "success create data",
	})
}

// @Summary Update Product
// @Description Update Product
// @Accept  json
// @Produce  json
// @Tags Product
// @Param id path string true "ID Product"
// @Router /api/product/{id} [patch]
// @Param payload body payload.UpdateProductRequest true "update product payload"
// @Success 200 {object} payload.Response200
// @Success 400 {object} payload.ResponseError
// @Success 404 {object} payload.ResponseError
// @Success 401 {object} payload.ResponseError
// @Success 500 {object} payload.ResponseError
func UpdateProduct(c *gin.Context) {
	var request payload.UpdateProductRequest

	tokenString := c.GetHeader("Authorization")
	tokenString = strings.ReplaceAll(tokenString, "Bearer ", "")
	claims, err := middleware.GetUserNameByToken(tokenString)
	if err != nil {
		c.JSON(http.StatusBadRequest, payload.ResponseError{
			Success: true,
			Message: "error",
			Error:   err.Error(),
		})
		return
	}

	err = json.NewDecoder(c.Request.Body).Decode(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, payload.ResponseError{
			Success: true,
			Message: "error",
			Error:   "invalid JSON",
		})
		return
	}

	validate := validator.New()
	// Validate struct
	err = validate.Struct(request)
	if err != nil {
		fmt.Println(err)
		errors := err.(validator.ValidationErrors)
		c.JSON(http.StatusBadRequest, payload.ResponseError{
			Success: true,
			Message: "error",
			Error:   fmt.Sprintf("%s", errors),
		})
		return
	}

	// mapping to Model
	var input models.UpdateProductInput
	_ = copier.Copy(&input, &request)

	// Get model if exist
	var product models.Product
	if err := database.Instance.Where("id = ?", c.Param("id")).First(&product).Error; err != nil {
		c.JSON(http.StatusBadRequest, payload.ResponseError{
			Success: true,
			Message: "error",
			Error:   "record not found",
		})
		return
	}

	input.UpdatedAt = time.Now()
	input.UpdatedBy = claims.Email

	database.Instance.Model(&product).Updates(input)

	c.JSON(http.StatusOK, payload.Response200{
		Success: true,
		Message: "success update data",
	})
}

// @Summary Update Product
// @Description Update Product
// @Accept  json
// @Produce  json
// @Tags Product
// @Param id path string true "ID Product"
// @Router /api/product/{id} [delete]
// @Success 200 {object} payload.Response200
// @Success 400 {object} payload.ResponseError
// @Success 404 {object} payload.ResponseError
// @Success 401 {object} payload.ResponseError
// @Success 500 {object} payload.ResponseError
func DeleteProduct(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	tokenString = strings.ReplaceAll(tokenString, "Bearer ", "")
	fmt.Println(tokenString)
	claims, err := middleware.GetUserNameByToken(tokenString)
	if err != nil {
		c.JSON(http.StatusBadRequest, payload.ResponseError{
			Success: true,
			Message: "error",
			Error:   err.Error(),
		})
		return
	}

	// Get model if exist
	var product models.Product
	if err := database.Instance.Where("id = ?", c.Param("id")).First(&product).Error; err != nil {
		c.JSON(http.StatusBadRequest,
			payload.ResponseError{
				Success: true,
				Message: "error",
				Error:   "record not found",
			})
		return
	}

	request := models.DeleteProductInput{
		UpdatedAt: time.Now(),
		DeletedAt: time.Now(),
		UpdatedBy: claims.Email,
		DeletedBy: claims.Email,
	}
	database.Instance.Model(&product).Updates(request)

	c.JSON(http.StatusOK, payload.Response200{
		Success: true,
		Message: "success deleted data",
	})
}
