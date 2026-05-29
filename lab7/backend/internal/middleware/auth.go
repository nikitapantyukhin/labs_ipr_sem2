package middleware

import (
	"fmt"
	"regexp"
	"sport_platform/internal/jwt"
	"sport_platform/internal/service_wrapper"

	"github.com/gin-gonic/gin"
)

const ClaimsKey = "Claims"

var authHeaderRegexp = regexp.MustCompile("Bearer (?P<token>\\S+)")
var tokenGroupIndex = authHeaderRegexp.SubexpIndex("token")

func AuthMiddleware(wrapper *service_wrapper.Wrapper) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header, exists := ctx.Request.Header["Authorization"]

		if !exists {
			fmt.Println("Authorization header is empty. Skipping...")
			ctx.Next()
			return
		}

		matches := authHeaderRegexp.FindStringSubmatch(header[0])
		if len(matches) <= tokenGroupIndex {
			fmt.Println("Token is of wrong format. Skipping...")
			ctx.Next()
			return
		}

		token := matches[tokenGroupIndex]
		claims, tokenValidationError := wrapper.JwtHandler.Validate(token, jwt.AccessToken)
		if tokenValidationError != nil {
			fmt.Printf("Error happend during token validation: %s\n", tokenValidationError)
			ctx.Next()
			return
		}

		ctx.Set(ClaimsKey, claims.Data)
		ctx.Next()
		return
	}
}
