package update_club_join_request

import (
	"sport_platform/application/models/shared"
	"time"
)

type UpdateClubJoinRequestResponse struct {
	ID        int64                    `json:"id"`
	ClubID    int64                    `json:"club_id"`
	UserID    int64                    `json:"user_id"`
	Status    shared.JoinRequestStatus `json:"status"`
	CreatedAt time.Time                `json:"created_at"`
	UpdatedAt time.Time                `json:"updated_at"`
}
