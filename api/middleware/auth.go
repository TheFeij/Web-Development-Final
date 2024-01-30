package middleware

import (
	"Messenger/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		accessToken := context.GetHeader("Authorization")
		claims, err := utils.ValidateToken(accessToken)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		context.Set("userClaims", claims)
		context.Next()
	}
}
