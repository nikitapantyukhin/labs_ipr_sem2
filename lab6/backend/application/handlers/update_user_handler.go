package handlers

import (
	"fmt"
	"net/http"
	"sport_platform/application/models/claims"
	"sport_platform/application/models/update_user"
	"sport_platform/internal/mapper"
	"sport_platform/internal/middleware"
	"sport_platform/internal/service_wrapper"
	"sport_platform/internal/sqlc/db_queries"

	"github.com/gin-gonic/gin"
)

func UpdateUserHandler(ctx *gin.Context, wrapper *service_wrapper.Wrapper) {
	var request update_user.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data"})
		return
	}

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

	var updateParams db_queries.UpdateUserParams

	var hashedPassword []byte
	if request.Password != nil {
		hashedPassword = wrapper.PasswordHandler.HashPassword(*request.Password)
	} else {
		hashedPassword = nil
	}

	paramsMappingError := mapper.Mapper{}.Map(
		&updateParams,
		request,
		struct {
			Password []byte
			ID       int64
		}{
			Password: hashedPassword,
			ID:       userClaims.ID,
		},
	)
	if paramsMappingError != nil {
		fmt.Println(paramsMappingError)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
		return
	}

	updatedUser, err := wrapper.Db.Queries.UpdateUser(ctx, updateParams)

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unknown error",
		})
		return
	}

	var response update_user.UpdateUserResponse

	responseMappingError := mapper.Mapper{}.Map(
		&response,
		updatedUser,
	)

	if responseMappingError != nil {
		fmt.Println(responseMappingError)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unknown error",
		})
		return
	}

	ctx.JSON(
		http.StatusOK,
		response,
	)
}
