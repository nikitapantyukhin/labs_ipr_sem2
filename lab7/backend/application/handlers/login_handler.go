package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"sport_platform/application/models/claims"
	"sport_platform/application/models/login"
	"sport_platform/internal/mapper"
	"sport_platform/internal/service_wrapper"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func LoginHandler(ctx *gin.Context, wrapper *service_wrapper.Wrapper) {
	var request login.LoginRequest
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"message": "can't decode request",
			},
		)
		return
	}

	user, dbError := wrapper.Db.Queries.GetUserByEmail(ctx, request.Email)
	switch {
	case errors.Is(dbError, pgx.ErrNoRows):
		ctx.JSON(
			http.StatusUnauthorized,
			gin.H{
				"message": "Email or password does not match",
			},
		)
		return
	case dbError != nil:
		fmt.Printf("Error happened during login: %s\n", dbError)
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "Something unusual happened",
			},
		)
		return
	}

	match, formatError := wrapper.PasswordHandler.VerifyPassword(request.Password, user.Password)
	if formatError != nil {
		fmt.Printf("Error happened during login: %s\n", formatError)
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "Something unusual happened",
			},
		)
		return
	}

	if !match {
		ctx.JSON(
			http.StatusUnauthorized,
			gin.H{
				"message": "Email or password does not match",
			},
		)
		return
	}

	var userClaims claims.UserClaims

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

	var response login.LoginResponse

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
		fmt.Printf("Error happened during login: %s\n", mappingError)
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
