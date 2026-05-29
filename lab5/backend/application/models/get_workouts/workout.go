package get_workouts

import "time"

type Workout struct {
	ID        int64     `json:"id"`
	ClubID    int64     `json:"club_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Cancelled bool      `json:"cancelled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
