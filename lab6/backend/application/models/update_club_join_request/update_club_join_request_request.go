package update_club_join_request

import "sport_platform/application/models/shared"

type UpdateClubJoinRequestRequest struct {
	ID     int64                    `json:"id"`
	Status shared.JoinRequestStatus `json:"status"`
}
