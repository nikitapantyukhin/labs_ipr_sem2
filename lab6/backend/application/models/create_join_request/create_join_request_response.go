package create_join_request

import (
	"sport_platform/application/models/shared"
)

type CreateJoinRequestResponse struct {
	ID     int64                    `json:"id"`
	ClubID int64                    `json:"club_id"`
	UserID int64                    `json:"user_id"`
	Status shared.JoinRequestStatus `json:"status"`
}
