package get_workouts

type GetWorkoutsRequest struct {
	ClubID int64 `uri:"id" binding:"required"`
}
