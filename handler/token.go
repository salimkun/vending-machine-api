package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/salimkun/vending-machine-api/database"
	"github.com/salimkun/vending-machine-api/middleware"
	"github.com/salimkun/vending-machine-api/models"
	"github.com/salimkun/vending-machine-api/payload"
)

// @Summary Generate Token User
// @Description Generate Token User
// @Accept  json
// @Produce  json
// @Tags Token
// @Router /api/token [post]
// @Param payload body payload.TokenRequest true "generate token payload"
// @Success 200 {object} payload.Response200
// @Success 400 {object} payload.ResponseError
// @Success 404 {object} payload.ResponseError
// @Success 401 {object} payload.ResponseError
// @Success 500 {object} payload.ResponseError
// @Security  ApiJwtToken
func GenerateToken(context *gin.Context) {
	var request payload.TokenRequest
	var user models.User

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	// check if email exists and password is correct
	record := database.Instance.Where("email = ?", request.Email).First(&user)
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}
	credentialError := user.CheckPassword(request.Password)
	if credentialError != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		context.Abort()
		return
	}
	tokenString, err := middleware.GenerateJWT(user.Email, user.Username)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	context.JSON(http.StatusOK, gin.H{"token": tokenString})
}
