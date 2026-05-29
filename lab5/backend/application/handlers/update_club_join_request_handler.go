package handlers

import (
	"fmt"
	"net/http"
	"sport_platform/application/models/claims"
	"sport_platform/application/models/shared"
	"sport_platform/application/models/update_club_join_request"
	"sport_platform/internal/mapper"
	"sport_platform/internal/middleware"
	"sport_platform/internal/service_wrapper"
	"sport_platform/internal/sqlc/db_queries"

	"github.com/gin-gonic/gin"
)

func UpdateClubJoinRequestStatusHandler(ctx *gin.Context, wrapper *service_wrapper.Wrapper) {
	var request update_club_join_request.UpdateClubJoinRequestRequest
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("can't parse query as error happend: %s", err),
		})
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

	if userClaims.Role != shared.Teacher {
		ctx.JSON(
			http.StatusForbidden,
			gin.H{
				"message": "No permission",
			},
		)
		return
	}

	switch request.Status {
	case shared.NotAccepted, shared.Accepted, shared.Declined:
	default:
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"message": "Invalid request status",
			},
		)
		return
	}

	var updateParams db_queries.UpdateClubJoinRequestStatusParams
	paramsMappingError := mapper.Mapper{}.Map(
		&updateParams,
		request,
	)

	if paramsMappingError != nil {
		fmt.Println(paramsMappingError)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
		return
	}

	updatedClubJoinRequest, err := wrapper.Db.Queries.UpdateClubJoinRequestStatus(ctx, updateParams)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unknown error",
		})
		return
	}

	var response update_club_join_request.UpdateClubJoinRequestResponse

	responseMappingError := mapper.Mapper{}.Map(
		&response,
		updatedClubJoinRequest,
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
