package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/salimkun/vending-machine-api/database"
	"github.com/salimkun/vending-machine-api/models"
	"github.com/salimkun/vending-machine-api/payload"
)

// @Summary Register User
// @Description Register User
// @Accept  json
// @Produce  json
// @Tags User
// @Router /api/user/register [post]
// @Param payload body payload.UserRequest true "create user payload"
// @Success 200 {object} payload.Response200
// @Success 400 {object} payload.ResponseError
// @Success 404 {object} payload.ResponseError
// @Success 401 {object} payload.ResponseError
// @Success 500 {object} payload.ResponseError
// @Security  ApiJwtToken
func RegisterUser(context *gin.Context) {
	var user models.User
	var userRequest payload.UserRequest
	if err := context.ShouldBindJSON(&userRequest); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	if err := user.HashPassword(userRequest.Password); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	_ = copier.Copy(user, userRequest)
	record := database.Instance.Create(&user)
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}
	context.JSON(http.StatusCreated, gin.H{"userId": user.ID, "email": user.Email, "username": user.Username})
}
