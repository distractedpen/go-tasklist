package auth

import (
	"errors"
	"go-tasklist/internal/user"
	"go-tasklist/internal/util"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(user.UserRequest) (user.UserDto, error)
	Register(user.UserRequest) error
	RefreshToken(user.UserDto) (user.UserDto, error)
}

type authService struct {
	userService user.UserService
}

func (a authService) Login(userRequest user.UserRequest) (user.UserDto, error) {
	// check if user exists
	existingUser, err := a.userService.GetUserByUsername(userRequest.Username)
	if nil != err {
		return user.UserDto{}, err
	}
	// check password correct
	if !validatePassword(existingUser.Password, userRequest.Password) {
		return user.UserDto{}, util.ErrAuthInvalid{}
	}

	// generate jwt token
	token, err := createToken(existingUser.Username)
	if nil != err {
		return user.UserDto{}, err
	}
	return user.UserDto{
		Id: existingUser.Id,
		Username: existingUser.Username,
		Token: token,
	}, nil
}

func (a authService) Register(userRequest user.UserRequest) error {
	// check if user exists
	_, err := a.userService.GetUserByUsername(userRequest.Username)
	if nil != err && !errors.Is(util.ErrDoesNotExist{}, err) {
		return err
	}

	// hash/salt password
	hashedPassword, err := generateHash(userRequest.Password)
	if nil != err {
		return err
	}
	// store new user
	err = a.userService.CreateUser(user.UserRequest{
		Username: userRequest.Username,
		Password: hashedPassword,
	})
	if nil != err {
		return err
	}
	return nil
}

func (a authService) RefreshToken(userDto user.UserDto) (user.UserDto, error) {
	newToken, err := createToken(userDto.Username)
	if nil != err {
		return userDto, err
	}

	return user.UserDto{
		Id: userDto.Id,
		Username: userDto.Username,
		Token: newToken,
	}, nil
}

func validatePassword(passwordHash string, password string) bool {
	log.Printf("%s, %s", passwordHash, password)
	err := bcrypt.CompareHashAndPassword(
		[]byte(passwordHash), []byte(password))
	return nil == err
}

func generateHash(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(password), bcrypt.DefaultCost)
	if nil != err {
		return "", err
	}
	return string(passwordHash), nil
}

func createToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})
	tokenstring, err := token.SignedString([]byte(os.Getenv("JWT_SIGNING_KEY")))
	if nil != err {
		return "", err
	}
	return tokenstring, nil
}

func VerifyToken(userToken string) (bool, error) {
	token, err := jwt.Parse(userToken, 
		func(t *jwt.Token) (any, error) {
			return []byte(os.Getenv("JWT_SIGNING_KEY")), nil
		})
	if nil != err {
		return false, err
	}

	if !token.Valid {
		return false, util.ErrAuthInvalid{}
	}

	return true, nil
}

func GetAuthService(userService user.UserService) AuthService {
	return authService{userService}
}
