package delete_workout

type DeleteWorkoutRequest struct {
	WorkoutID int64 `uri:"workout_id" binding:"required"`
}
