package get_join_request

type GetJoinRequestRequest struct {
	ID     *int64 `form:"id"`
	ClubID *int64 `form:"club_id"`
	UserID *int64 `form:"user_id"`
	Limit  *int64 `form:"limit_"`
	Offset *int64 `form:"offset_"`
}
