package middleware

import (
	"Messenger/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		accessToken := context.GetHeader("Authorization")

		if accessToken == "" {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - Missing access token"})
			return
		}

		parsedToken := utils.ParseToken(accessToken)

		if parsedToken == nil || !parsedToken.Valid {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - Invalid access token"})
			return
		}

		// Attach the user claims to the context for use in the handlers
		claims, ok := parsedToken.Claims.(*utils.UserClaims)
		if !ok {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - Invalid token claims"})
			return
		}

		context.Set("userClaims", claims)
		context.Next()
	}
}
