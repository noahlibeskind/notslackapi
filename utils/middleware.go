// NotSlack API
// created by Noah Libeskind
// the functions used here are closely derived from https://seefnasrul.medium.com/create-your-first-go-rest-api-with-jwt-authentication-in-gin-framework-dbe5bda72817

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
