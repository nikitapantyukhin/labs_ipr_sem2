package handlers

import (
	"fmt"
	"net/http"
	"sport_platform/application/models/get_clubs"
	"sport_platform/internal/mapper"
	"sport_platform/internal/middleware"
	"sport_platform/internal/service_wrapper"

	"github.com/gin-gonic/gin"
)

func GetClubHandler(ctx *gin.Context, wrapper *service_wrapper.Wrapper) {
	var request get_clubs.GetClubRequest
	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("can't parse query as error happend: %s", err),
		})
		return
	}
	_, exists := ctx.Get(middleware.ClaimsKey)
	if !exists {
		ctx.JSON(
			http.StatusUnauthorized,
			gin.H{
				"message": "Unauthorized",
			},
		)
		return
	}

	club, dbError := wrapper.Db.Queries.GetClubById(ctx, request.ID)
	if dbError != nil {
		fmt.Printf("Error happened getting club: %s\n", dbError)

		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "Something unusual happened",
			},
		)
		return
	}
	clubAttachments, dbError := wrapper.Db.Queries.GetClubAttachments(ctx, request.ID)
	if dbError != nil {
		fmt.Printf("Error happened getting club attachments: %s\n", dbError)

		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "Something unusual happened",
			},
		)
		return
	}

	var response get_clubs.GetClubResponse
	mappingError := mapper.Mapper{}.Map(
		&response,
		club,
		struct {
			Attachments []string
		}{
			Attachments: clubAttachments,
		},
	)

	if mappingError != nil {
		fmt.Printf("Club mapping error: %s\n", mappingError)
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
