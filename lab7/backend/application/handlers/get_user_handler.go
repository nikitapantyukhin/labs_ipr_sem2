package handlers

import (
	"fmt"
	"net/http"
	"sport_platform/application/models/claims"
	"sport_platform/application/models/get_user"
	"sport_platform/internal/mapper"
	"sport_platform/internal/middleware"
	"sport_platform/internal/service_wrapper"

	"github.com/gin-gonic/gin"
)

func GetUserHandler(ctx *gin.Context, wrapper *service_wrapper.Wrapper) {
	claimsRaw, exists := ctx.Get(middleware.ClaimsKey)
	if !exists {
		ctx.JSON(
			http.StatusUnauthorized,
			gin.H{
				"message": "Unauthorized",
			},
		)
		return
	}

	userClaims := claimsRaw.(claims.UserClaims)
	user, dbError := wrapper.Db.Queries.GetUserById(ctx, userClaims.ID)
	if dbError != nil {
		fmt.Printf("Error while getting user: %s\n", dbError)
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "Something unusual happened",
			},
		)
		return
	}

	claimsMappingError := mapper.Mapper{}.Map(&userClaims, user)
	if claimsMappingError != nil {
		fmt.Println(claimsMappingError)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unknown error",
		})
		return
	}

	accessToken, refreshToken, tokenGenerationError := wrapper.JwtHandler.GenerateJwtPair(userClaims, fmt.Sprintf("%d", user.ID))
	if tokenGenerationError != nil {
		fmt.Println(claimsMappingError)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unknown error",
		})
		return
	}

	var response get_user.GetUserResponse

	mappingError := mapper.Mapper{}.Map(
		&response,
		user,
		struct {
			AccessToken  string
			RefreshToken string
		}{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	)
	if mappingError != nil {
		fmt.Printf("Error happened while getting user: %s\n", mappingError)
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "Something unusual happened",
			},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		response,
	)
}
