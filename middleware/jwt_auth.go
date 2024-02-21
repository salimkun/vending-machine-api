package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/salimkun/vending-machine-api/payload"
)

func AuthJwt() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenString := context.GetHeader("Authorization")
		tokenString = strings.ReplaceAll(tokenString, "Bearer ", "")

		if tokenString == "" {
			context.JSON(401, payload.ResponseError{
				Success: false,
				Message: "error",
				Error:   "request does not contain an access token",
			})
			context.Abort()
			return
		}
		err := ValidateToken(tokenString)
		if err != nil {
			context.JSON(401, payload.ResponseError{
				Success: false,
				Message: "error",
				Error:   err.Error(),
			})
			context.Abort()
			return
		}
		context.Next()
	}
}
