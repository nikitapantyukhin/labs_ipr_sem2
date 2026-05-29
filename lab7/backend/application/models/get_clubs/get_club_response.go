package get_clubs

import "time"

type GetClubResponse struct {
	ID                     int64     `json:"id"`
	Name                   string    `json:"name"`
	Description            string    `json:"description"`
	SportTypeID            int64     `json:"sport_type_id"`
	SportTypeName          string    `json:"sport_type_name"`
	TeacherID              int64     `json:"teacher_id"`
	TeacherName            string    `json:"teacher_name"`
	TotalPlaces            int32     `json:"total_places,omitempty"`
	Place                  string    `json:"place"`
	EducationLevelID       int64     `json:"education_level_id"`
	EducationLevelName     string    `json:"education_level_name"`
	RequiredWorkoutPerWeek int32     `json:"required_workout_per_week"`
	Attachments            []string  `json:"attachments,omitempty"`
	CreatedAt              time.Time `json:"created_at"`
}

type GetAllClubsResponse struct {
	Clubs []GetClubResponse `json:"clubs"`
}
