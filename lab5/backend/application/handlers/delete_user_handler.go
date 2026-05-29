package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"sport_platform/application/models/claims"
	"sport_platform/application/models/delete_user"
	"sport_platform/application/models/shared"
	"sport_platform/internal/middleware"
	"sport_platform/internal/service_wrapper"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func DeleteUserHandler(ctx *gin.Context, wrapper *service_wrapper.Wrapper) {
	var request delete_user.DeleteUserRequest
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

	if userClaims.Email != request.Email && userClaims.Role != shared.Admin {
		ctx.JSON(
			http.StatusForbidden,
			gin.H{
				"message": "No permission",
			},
		)
		return
	}

	user, dbError := wrapper.Db.Queries.GetUserByEmail(ctx, request.Email)
	switch {
	case errors.Is(dbError, pgx.ErrNoRows):
		ctx.JSON(
			http.StatusNotFound,
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
		fmt.Printf("Error happened during password verification: %s\n", formatError)
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

	if err := wrapper.Db.Queries.DeleteUser(ctx, request.Email); err != nil {
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
		gin.H{
			"message": "User deleted",
		},
	)
}
