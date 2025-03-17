package user

import (
	"database/sql"
	"errors"
	"log"

	"go-tasklist/internal/util"
)


type UserRepository interface {
	GetUserById(id int64) (User, error)
	GetUserByUsername(username string) (User, error)
	CreateUser(UserRequest) error
}

type userRepository struct {
	db *sql.DB
}

func (r userRepository) GetUserById(id int64) (User, error) {
	var user User

	res := r.db.QueryRow(
		`SELECT * FROM users WHERE id=?`,
		id);

	err := res.Scan(&user.Id, &user.Username, &user.Password)
	if sql.ErrNoRows == err {
		return user, util.ErrDoesNotExist{}
	}
	if nil != err && sql.ErrNoRows != err {
		log.Println(err.Error())
		return user, util.ErrDB{}
	}

	return user, nil
}

func (r userRepository) GetUserByUsername(username string) (User, error) {
	var user User

	res := r.db.QueryRow(
		`SELECT * FROM users WHERE username=?`,
		username);

	err := res.Scan(&user.Id, &user.Username, &user.Password)
	if sql.ErrNoRows == err {
		return user, util.ErrDoesNotExist{}
	}
	if nil != err && sql.ErrNoRows != err {
		log.Println(err.Error())
		return user, util.ErrDB{}
	}

	return user, nil
}

func (r userRepository) CreateUser(userRequest UserRequest) error {
	_, err := r.GetUserByUsername(userRequest.Username)
	if errors.Is(util.ErrDoesNotExist{}, err) {
		_, err = r.db.Exec("INSERT INTO users(username, password) VALUES (?, ?)", 
			userRequest.Username, userRequest.Password);
		if nil != err {
			log.Println(err)
			return util.ErrDB{}
		}
		return nil

	} else if nil != err && sql.ErrNoRows != err {
		log.Println(err)
		return util.ErrDB{}
	} else {
		log.Println(err)
		return util.ErrExists{}
	}
} 

func GetUserRepository(db *sql.DB) UserRepository {
	return userRepository{db}
}
