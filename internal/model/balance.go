package model

type Balance struct {
	ID        int
	UserID    int
	Current   float64
	Withdrawn float64
}
