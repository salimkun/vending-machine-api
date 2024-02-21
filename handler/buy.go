package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/salimkun/vending-machine-api/database"
	"github.com/salimkun/vending-machine-api/models"
	"github.com/salimkun/vending-machine-api/payload"
)

// @Summary Buy Product
// @Description Buy Product
// @Accept  json
// @Produce  json
// @Tags Product
// @Router /api/buy/product [post]
// @Param payload body payload.BuyProductRequest true "update buy product payload"
// @Success 200 {object} payload.Response200
// @Success 400 {object} payload.ResponseError
// @Success 404 {object} payload.ResponseError
// @Success 401 {object} payload.ResponseError
// @Success 500 {object} payload.ResponseError
func BuyProduct(c *gin.Context) {
	var request payload.BuyProductRequest

	err := json.NewDecoder(c.Request.Body).Decode(&request)
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

	var products []models.Product
	database.Instance.Order("price desc").Find(&products)

	result, err := vendingMachine(request.Money, products)
	if err != nil {
		c.JSON(http.StatusBadRequest, payload.ResponseError{
			Success: true,
			Message: "error",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, payload.Response200{
		Success: true,
		Message: "success update data",
		Data:    createKeyValuePairs(result),
	})
}

func vendingMachine(money []int, stock []models.Product) (map[string]int, error) {
	total := 0
	product := make(map[string]int)

	for _, i := range money {
		switch i {
		case 2000, 5000:
			total = total + i
		default:
			return nil, errors.New("invalid denomination")
		}

	}

	minPrice := 0
	for total > 0 {
		for _, i := range stock {
			if total-i.Price >= 0 {
				product[i.Name] += 1
				total = total - i.Price
				minPrice = i.Price
			}
		}

		if total < minPrice {
			break
		}
	}

	return product, nil

}

func createKeyValuePairs(m map[string]int) string {
	b := new(bytes.Buffer)
	count := 0
	for key, value := range m {
		count++
		if count < len(m) {
			fmt.Fprintf(b, "%d %s, ", value, key)
		} else {
			fmt.Fprintf(b, "%d %s ", value, key)
		}
	}
	return b.String()
}
