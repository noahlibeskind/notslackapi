package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		err := TokenValid(context)
		if err != nil {
			context.String(http.StatusUnauthorized, "Unauthorized")
			context.Abort()
			return
		}
		context.Next()
	}
}
