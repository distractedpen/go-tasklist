package user

import (
	"database/sql"
	"fmt"
	"log"
)

type UserRepository interface {
	GetUser(int64)
}

type userRepository struct {
	db *sql.DB
}

func (r *userRepository) GetUser(userId int64) (UserDto, error) {
	var user UserDto

	res := r.db.QueryRow(
		`SELECT id, username FROM users WHERE userId=?`,
		userId);

	err := res.Scan(&user.Id, &user.Username)
	if nil != err {
		log.Println(err.Error())
		return user, fmt.Errorf("Error Reading from DB.")
	}

	return user, nil
}

