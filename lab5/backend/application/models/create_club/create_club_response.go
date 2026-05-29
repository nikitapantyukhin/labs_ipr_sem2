package create_club

import "time"

type CreateClubResponse struct {
	ID                     int64     `json:"id"`
	Name                   string    `json:"name"`
	Description            string    `json:"description"`
	SportTypeID            int64     `json:"sport_type_id"`
	TeacherID              int64     `json:"teacher_id"`
	TotalPlaces            int       `json:"total_places"`
	Place                  string    `json:"place"`
	EducationLevelID       int64     `json:"education_level_id"`
	RequiredWorkoutPerWeek int       `json:"required_workout_per_week"`
	Attachments            []string  `json:"attachments"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}
