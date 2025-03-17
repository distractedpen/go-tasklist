package user

type User struct {
	Id int64
	Username string
	Password string
}

type UserDto struct {
	Id int64             `json:"id"`
	Username string      `json:"username"`
	Token string         `json:"token"`
}

type UserRequest struct {
	Username string      `json:"username"`
	Password string      `json:"password"`
}


func (u User) MapToUserDto(token string) UserDto {
	return UserDto{
		Id: u.Id,
		Username: u.Username,
		Token: token,
	}
}
