package create_club

type CreateClubRequest struct {
	Name                   string `json:"name" form:"name" validate:"required,min=3,max=100"`
	Description            string `json:"description" form:"description" validate:"required"`
	SportTypeID            int64  `json:"sport_type_id" form:"sport_type_id" validate:"required,gte=1"`
	TeacherID              int64  `json:"teacher_id" form:"teacher_id"`
	TotalPlaces            int    `json:"total_places" form:"total_places" validate:"gte=0"`
	Place                  string `json:"place" form:"place" validate:"required"`
	EducationLevelName     string `json:"education_level_name" form:"education_level_name" validate:"required,min=2"`
	RequiredWorkoutPerWeek int    `json:"required_workout_per_week" form:"required_workout_per_week" validate:"required,gte=1"`
}
