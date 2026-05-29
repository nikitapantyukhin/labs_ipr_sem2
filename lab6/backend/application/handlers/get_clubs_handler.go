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

func GetClubsHandler(ctx *gin.Context, wrapper *service_wrapper.Wrapper) {
	var request get_clubs.GetClubRequest
	if err := ctx.ShouldBind(&request); err != nil {
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

	clubs, dbError := wrapper.Db.Queries.GetAllClubs(ctx)
	if dbError != nil {
		fmt.Printf("Error while getting clubs: %s\n", dbError)

		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "Something unusual happened",
			},
		)
		return
	}

	response := get_clubs.GetAllClubsResponse{
		Clubs: make([]get_clubs.GetClubResponse, len(clubs)),
	}

	for idx, club := range clubs {
		attachments, dbError := wrapper.Db.Queries.GetClubAttachments(ctx, club.ID)
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

		mappingError := mapper.Mapper{}.Map(
			&response.Clubs[idx],
			club,
			struct {
				Attachments []string
			}{
				Attachments: attachments,
			},
		)
		if mappingError != nil {
			fmt.Printf("Clubs mapping error: %s\n", mappingError)
			ctx.JSON(
				http.StatusInternalServerError,
				gin.H{
					"message": "Unknown error",
				},
			)
			return
		}
	}

	ctx.JSON(
		http.StatusOK,
		response,
	)
}
