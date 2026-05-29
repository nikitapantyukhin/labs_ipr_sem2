package get_join_request

import (
	"sport_platform/application/models/shared"
	"time"
)

type GetJoinRequestResponse struct {
	ID        int64                    `json:"id"`
	ClubID    int64                    `json:"club_id"`
	UserID    int64                    `json:"user_id"`
	Status    shared.JoinRequestStatus `json:"status"`
	CreatedAt time.Time                `json:"created_at"`
	UpdatedAt time.Time                `json:"updated_at"`
}

type JoinRequest struct {
	ID        int64                    `json:"id"`
	ClubID    int64                    `json:"club_id"`
	UserID    int64                    `json:"user_id"`
	Status    shared.JoinRequestStatus `json:"status,omitempty"`
	CreatedAt time.Time                `json:"created_at"`
	UpdatedAt time.Time                `json:"updated_at"`
}

type GetJoinRequestsResponse struct {
	JoinRequests []JoinRequest `json:"join_requests"`
}
