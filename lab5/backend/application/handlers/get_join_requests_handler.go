package handlers

import (
	"fmt"
	"net/http"
	"sport_platform/application/models/claims"
	"sport_platform/application/models/get_join_request"
	"sport_platform/internal/mapper"
	"sport_platform/internal/middleware"
	"sport_platform/internal/service_wrapper"
	"sport_platform/internal/sqlc/db_queries"

	"github.com/gin-gonic/gin"
)

func GetJoinRequestsHandler(ctx *gin.Context, wrapper *service_wrapper.Wrapper) {
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
	var request get_join_request.GetJoinRequestRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		fmt.Printf("GetJoinRequestsHandler: ShouldBind error: %s\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Unknown error",
		})
		return
	}

	var getParams db_queries.GetJoinRequestsParams

	paramsMappingError := mapper.Mapper{}.Map(
		&getParams,
		request,
		struct {
			CurrentUserID int64
		}{
			CurrentUserID: user.ID,
		},
	)
	if paramsMappingError != nil {
		fmt.Println(paramsMappingError)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unknown error",
		})
		return
	}
	joinRequests, dbError := wrapper.Db.Queries.GetJoinRequests(ctx, getParams)
	if dbError != nil {
		fmt.Printf("Error while getting join requests: %s\n", dbError)

		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "Something unusual happened",
			},
		)
		return
	}

	var response get_join_request.GetJoinRequestsResponse
	response.JoinRequests = make([]get_join_request.JoinRequest, len(joinRequests))

	type SourceWrapper struct {
		Items []db_queries.ClubJoinRequest
	}

	type DestWrapper struct {
		Items []get_join_request.JoinRequest
	}

	sourceWrapper := SourceWrapper{Items: joinRequests}
	destWrapper := DestWrapper{Items: response.JoinRequests}

	mappingError := mapper.Mapper{}.Map(&destWrapper, sourceWrapper)
	if mappingError != nil {
		fmt.Printf("GetAllJoinRequestsHandler: mapping error: %s\n", mappingError)
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "Unknown error",
			},
		)
		return
	}
	response.JoinRequests = destWrapper.Items

	ctx.JSON(
		http.StatusOK,
		response,
	)
}
