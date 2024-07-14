package model

type User struct {
	Id       int    `db:"id"`
	Login    string `db:"login"`
	Password string `db:"password"`
}
