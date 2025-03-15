package user

type User struct {
	Id int64
	Username string
	Password string
}

type UserDto struct {
	Id int64             `json:"id"`
	Username string      `json:"username"`
}
