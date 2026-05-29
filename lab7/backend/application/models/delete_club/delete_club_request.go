package delete_club

type DeleteClubRequest struct {
	ID int64 `uri:"id" binding:"required"`
}
