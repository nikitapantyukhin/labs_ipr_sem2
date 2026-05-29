package handlers

import (
	"fmt"
	"net/http"
	"sport_platform/application/models/get_workouts"
	"sport_platform/internal/mapper"
	"sport_platform/internal/service_wrapper"
	"sport_platform/internal/sqlc/db_queries"

	"github.com/gin-gonic/gin"
)

type WorkoutsDBWrapper struct {
	Workouts []db_queries.GetWorkoutsByClubRow
}

type WorkoutsResponseWrapper struct {
	Workouts []get_workouts.Workout
}

func GetClubWorkoutsHandler(ctx *gin.Context, wrapper *service_wrapper.Wrapper) {
	var err error
	var request get_workouts.GetWorkoutsRequest
	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid club ID format in URL"})
		return
	}

	existsClub, err := wrapper.Db.Queries.CheckClubExists(ctx, request.ClubID)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Database error while checking club existence"})
		return
	}

	if !existsClub {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Club not found"})
		return
	}

	workoutsDB, err := wrapper.Db.Queries.GetWorkoutsByClub(ctx, request.ClubID)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Database error retrieving workouts"})
		return
	}

	dbWrapper := WorkoutsDBWrapper{Workouts: workoutsDB}
	var responseWrapper WorkoutsResponseWrapper
	err = mapper.Mapper{}.Map(&responseWrapper, dbWrapper)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
		return
	}

	response := get_workouts.GetWorkoutsResponse{
		Workout: responseWrapper.Workouts,
	}

	ctx.JSON(http.StatusOK, response)
}
