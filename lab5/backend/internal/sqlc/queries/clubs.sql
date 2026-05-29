-- name: CreateClub :one
INSERT INTO clubs (
    name,
    description,
    sport_type_id,
    teacher_id,
    total_places,
    place,
    education_level_id,
    required_workout_per_week
) VALUES (
             sqlc.arg(name),
             sqlc.arg(description),
             sqlc.arg(sport_type_id),
             sqlc.arg(teacher_id),
             sqlc.arg(total_places),
             sqlc.arg(place),
             (SELECT id FROM education_levels WHERE education_levels.name = sqlc.arg(education_level_name)),
             sqlc.arg(required_workout_per_week)
         )
    RETURNING *;

-- name: GetClubByID :one
SELECT * FROM clubs WHERE id = $1 AND is_deleted = false;

-- name: SoftDeleteClub :exec
UPDATE clubs
SET is_deleted = true, updated_at = NOW()
WHERE id = $1;

-- name: GetClubsByTeacher :many
SELECT * FROM clubs
WHERE teacher_id = $1 AND is_deleted = false
ORDER BY created_at DESC;

-- name: ListActiveClubs :many
SELECT * FROM clubs
WHERE is_deleted = false
ORDER BY name;

-- name: CheckClubOwnership :one
SELECT EXISTS(
    SELECT 1 FROM clubs
    WHERE id = $1 AND teacher_id = $2 AND is_deleted = false
);
-- name: CheckClubExists :one
SELECT EXISTS(
    SELECT 1 FROM clubs
    WHERE id = $1 AND is_deleted = FALSE
);

-- name: CheckSportTypeExists :one
SELECT EXISTS(
    SELECT 1 FROM sport_types
    WHERE id = $1
);

-- name: CheckEducationLevelExistsByName :one
SELECT EXISTS(
    SELECT 1 FROM education_levels
    WHERE name = $1
);

-- name: GetAllClubs :many
SELECT
    clubs.id,
    clubs.name,
    clubs.description,
    clubs.sport_type_id,
    clubs.teacher_id,
    clubs.total_places,
    clubs.place,
    clubs.education_level_id,
    clubs.required_workout_per_week,
    clubs.created_at,
    clubs.updated_at,
    clubs.is_deleted,
    sport_types.name as sport_type_name,
    education_levels.name as education_level_name,
    users.full_name as teacher_name
FROM clubs
        JOIN sport_types ON clubs.sport_type_id = sport_types.id AND sport_types.is_deleted = false
        JOIN education_levels ON clubs.education_level_id = education_levels.id AND education_levels.is_deleted = false
        JOIN users ON clubs.teacher_id = users.id AND users.is_deleted = false
WHERE clubs.is_deleted = false;

-- name: GetClubById :one
SELECT
    clubs.id,
    clubs.name,
    clubs.description,
    clubs.sport_type_id,
    clubs.teacher_id,
    clubs.total_places,
    clubs.place,
    clubs.education_level_id,
    clubs.required_workout_per_week,
    clubs.created_at,
    clubs.updated_at,
    clubs.is_deleted,
    sport_types.name as sport_type_name,
    education_levels.name as education_level_name,
    users.full_name as teacher_name
FROM clubs
        JOIN sport_types ON clubs.sport_type_id = sport_types.id AND sport_types.is_deleted = false
        JOIN education_levels ON clubs.education_level_id = education_levels.id AND education_levels.is_deleted = false
        JOIN users ON clubs.teacher_id = users.id AND users.is_deleted = false
WHERE clubs.id = @id AND clubs.is_deleted = false;

-- name: GetClubAttachments :many
SELECT attachment_url
FROM club_attachments
WHERE club_id = @id;

-- name: UploadAttachment :exec
INSERT INTO club_attachments
    (club_id, attachment_url)
VALUES
    (@club_id, @attachment_url);
