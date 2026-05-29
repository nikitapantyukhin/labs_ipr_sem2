-- name: GetWorkoutsByClub :many
SELECT id, club_id, start_date, end_date, cancelled, created_at, updated_at
FROM workouts
WHERE club_id = $1 AND is_deleted = FALSE AND cancelled = FALSE
ORDER BY start_date DESC;

-- name: CheckWorkoutExists :one
SELECT EXISTS(
    SELECT 1 FROM workouts
    WHERE id = $1 AND is_deleted = FALSE
);

-- name: SoftDeleteWorkout :exec
UPDATE workouts
SET is_deleted = TRUE, updated_at = NOW()
WHERE id = $1;

-- name: CheckWorkoutOwnership :one
SELECT EXISTS (
    SELECT 1 FROM workouts w
                      JOIN clubs c ON w.club_id = c.id
    WHERE w.id = $1        
      AND c.teacher_id = $2
      AND w.is_deleted = FALSE
);

-- name: CreateWorkout :one
INSERT INTO workouts
(club_id, start_date, end_date)
VALUES
(@club_id, @start_date, @end_date)
RETURNING workouts.*;

-- name: UpdateWorkout :one
UPDATE workouts
SET
    cancelled = COALESCE(sqlc.narg(cancelled), cancelled),
    start_date = COALESCE(sqlc.narg(start_date), start_date),
    end_date = COALESCE(sqlc.narg(end_date), end_date),
    updated_at = now()
WHERE id = @id
RETURNING workouts.*;
