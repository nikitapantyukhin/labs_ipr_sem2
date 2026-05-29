package application

import (
	"context"
	"time"
)

const demoPassword = "Demo12345!"

func (appl *Application) SeedDemoData(ctx context.Context) error {
	queries := []struct {
		sql  string
		args []interface{}
	}{
		{
			sql: "INSERT INTO sport_types (name) VALUES ($1), ($2), ($3), ($4) ON CONFLICT DO NOTHING",
			args: []interface{}{
				"Футбол",
				"Баскетбол",
				"Волейбол",
				"Легкая атлетика",
			},
		},
		{
			sql: "INSERT INTO education_levels (name) VALUES ($1), ($2), ($3) ON CONFLICT DO NOTHING",
			args: []interface{}{
				"Начальный",
				"Средний",
				"Продвинутый",
			},
		},
		{
			sql: "INSERT INTO users (full_name, social_network_link, phone_number, email, birth_date, role, password) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT DO NOTHING",
			args: []interface{}{
				"Алексей Морозов",
				"@coach_morozov",
				"+79990000001",
				"coach@sport.local",
				time.Date(1990, time.May, 20, 0, 0, 0, 0, time.UTC),
				"Teacher",
				appl.wrapper.PasswordHandler.HashPassword(demoPassword),
			},
		},
		{
			sql: `
INSERT INTO clubs (name, description, sport_type_id, teacher_id, total_places, place, education_level_id, required_workout_per_week)
SELECT $1, $2, sport_types.id, users.id, $3, $4, education_levels.id, $5
FROM sport_types, users, education_levels
WHERE sport_types.name = $6
  AND users.email = $7
  AND education_levels.name = $8
ON CONFLICT DO NOTHING`,
			args: []interface{}{
				"Футбол для начинающих",
				"Тренировки по технике, координации и командной игре для студентов без опыта.",
				24,
				"Стадион кампуса",
				2,
				"Футбол",
				"coach@sport.local",
				"Начальный",
			},
		},
		{
			sql: `
INSERT INTO clubs (name, description, sport_type_id, teacher_id, total_places, place, education_level_id, required_workout_per_week)
SELECT $1, $2, sport_types.id, users.id, $3, $4, education_levels.id, $5
FROM sport_types, users, education_levels
WHERE sport_types.name = $6
  AND users.email = $7
  AND education_levels.name = $8
ON CONFLICT DO NOTHING`,
			args: []interface{}{
				"Баскетбольная секция",
				"Регулярные тренировки, игровые разборы и подготовка к университетским турнирам.",
				18,
				"Спортзал N2",
				3,
				"Баскетбол",
				"coach@sport.local",
				"Средний",
			},
		},
		{
			sql: `
INSERT INTO clubs (name, description, sport_type_id, teacher_id, total_places, place, education_level_id, required_workout_per_week)
SELECT $1, $2, sport_types.id, users.id, $3, $4, education_levels.id, $5
FROM sport_types, users, education_levels
WHERE sport_types.name = $6
  AND users.email = $7
  AND education_levels.name = $8
ON CONFLICT DO NOTHING`,
			args: []interface{}{
				"Волейбол микс",
				"Секция для тех, кто хочет играть стабильно, расти по уровню и участвовать в матчах.",
				20,
				"Большой зал",
				2,
				"Волейбол",
				"coach@sport.local",
				"Начальный",
			},
		},
	}

	for _, query := range queries {
		if err := appl.wrapper.Db.Exec(ctx, query.sql, query.args...); err != nil {
			return err
		}
	}

	workouts := []struct {
		clubName  string
		startDate time.Time
		endDate   time.Time
	}{
		{
			clubName:  "Футбол для начинающих",
			startDate: time.Date(2030, time.January, 15, 18, 0, 0, 0, time.UTC),
			endDate:   time.Date(2030, time.January, 15, 19, 30, 0, 0, time.UTC),
		},
		{
			clubName:  "Баскетбольная секция",
			startDate: time.Date(2030, time.January, 16, 17, 30, 0, 0, time.UTC),
			endDate:   time.Date(2030, time.January, 16, 19, 0, 0, 0, time.UTC),
		},
		{
			clubName:  "Волейбол микс",
			startDate: time.Date(2030, time.January, 17, 18, 30, 0, 0, time.UTC),
			endDate:   time.Date(2030, time.January, 17, 20, 0, 0, 0, time.UTC),
		},
	}

	for _, workout := range workouts {
		if err := appl.wrapper.Db.Exec(
			ctx,
			`
INSERT INTO workouts (club_id, start_date, end_date)
SELECT clubs.id, $2, $3
FROM clubs
WHERE clubs.name = $1
  AND NOT EXISTS (
    SELECT 1
    FROM workouts
    WHERE workouts.club_id = clubs.id
      AND workouts.start_date = $2
      AND workouts.is_deleted = false
  )`,
			workout.clubName,
			workout.startDate,
			workout.endDate,
		); err != nil {
			return err
		}
	}

	return nil
}
