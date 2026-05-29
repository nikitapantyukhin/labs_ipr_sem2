package handlers

import (
	"fmt"
	"net/http"
	"sport_platform/application/models/claims"
	"sport_platform/application/models/delete_workout"
	"sport_platform/application/models/shared"
	"sport_platform/internal/middleware"
	"sport_platform/internal/service_wrapper"
	"sport_platform/internal/sqlc/db_queries"

	"github.com/gin-gonic/gin"
)

func DeleteWorkoutHandler(ctx *gin.Context, wrapper *service_wrapper.Wrapper) {
	claimsRaw, exists := ctx.Get(middleware.ClaimsKey)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	userClaims := claimsRaw.(claims.UserClaims)
	if userClaims.Role != shared.Admin && userClaims.Role != shared.Teacher {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "No permission"})
		return
	}

	var request delete_workout.DeleteWorkoutRequest
	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid workout ID format in URL"})
		return
	}

	exists, err := wrapper.Db.Queries.CheckWorkoutExists(ctx, request.WorkoutID)

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Database error while checking workout existence"})
		return
	}

	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Workout not found or already deleted"})
		return
	}

	if userClaims.Role != shared.Admin {

		var isOwner bool
		var err error
		isOwner, err = wrapper.Db.Queries.CheckWorkoutOwnership(ctx, db_queries.CheckWorkoutOwnershipParams{
			ID:        request.WorkoutID,
			TeacherID: userClaims.ID,
		})

		if err != nil {
			fmt.Println(err)
			ctx.JSON(
				http.StatusInternalServerError,
				gin.H{"message": "Database error while checking ownership"},
			)
			return
		}

		if !isOwner {
			ctx.JSON(
				http.StatusForbidden,
				gin.H{"message": "You can only delete workouts belonging to your club"},
			)
			return
		}
	}

	err = wrapper.Db.Queries.SoftDeleteWorkout(ctx, request.WorkoutID)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"workout_id": request.WorkoutID,
		"message":    "Workout deleted successfully",
	})
}
