package create_workout

import (
	"time"
)

type CreateWorkoutRequest struct {
	ClubID    int64     `json:"club_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}
