package domain

type TripModel struct {
	ID       string
	UserID   string
	Status   string
	RideFare *RideFareModel
}
