package auth

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"go-tasklist/internal/user"
	"go-tasklist/internal/util"
)


type AuthAPI interface {
	Login(http.ResponseWriter, *http.Request)
	Register(http.ResponseWriter, *http.Request)
}

type authApi struct {
	authService AuthService
}


func (a authApi) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userRequest, err := parseRequestBody(r)
	if nil != err {
		log.Println(err)
		util.SendResponse(w, 400, nil)
		return
	}
	log.Printf("%v", userRequest)
	userDto, err := a.authService.Login(userRequest)
	if nil != err {
		log.Println(err)
		util.SendResponse(w, 500, nil)
		return
	}

	util.SendResponse(w, 200, userDto)
}

func (a authApi) Register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	response := make(map[string]any)

	userRequest, err := parseRequestBody(r)
	if nil != err {
		log.Println(err)
		util.SendResponse(w, 400, nil)
		return
	}

	err = a.authService.Register(userRequest)
	if nil != err {
		log.Printf("Error Creating User %v\n", err)
		util.SendResponse(w, 500, nil)
		return
	}
	util.SendResponse(w, 200, response)
}


func GetAuthHandlers(authService AuthService) AuthAPI {
	return authApi{
		authService: authService,
	};
}

func parseRequestBody(r *http.Request) (user.UserRequest, error){
	var userRequest user.UserRequest
	bodyBytes, err := io.ReadAll(r.Body)
	if nil != err {
		return userRequest, util.ErrRequestInvalid{}
	}

	err = json.Unmarshal(bodyBytes, &userRequest)
	if nil != err {
		return userRequest, util.ErrRequestInvalid{}
	}
	return userRequest, nil
}
