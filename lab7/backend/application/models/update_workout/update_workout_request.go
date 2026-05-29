package update_workout

import "time"

type UpdateWorkoutRequest struct {
	ID        int64      `json:"id"`
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	Cancelled *bool      `json:"cancelled"`
}
