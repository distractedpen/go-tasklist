package user

type UserService interface {
	GetUserByUsername(username string) (User, error)
	GetUserById(id int64) (User, error)
	CreateUser(UserRequest) error
}

type userService struct {
	userRepository UserRepository
}


func (s userService) GetUserById(id int64) (User, error) {
	user, err := s.userRepository.GetUserById(id)
	return user, err
}

func (s userService) GetUserByUsername(username string) (User, error) {
	user, err := s.userRepository.GetUserByUsername(username)
	return user, err
}

func (s userService) CreateUser(userRequset UserRequest) error {
	err := s.userRepository.CreateUser(userRequset)
	if nil != err {
		return err
	}
	return nil
}

func GetUserService(userRepository UserRepository) UserService {
	return userService{
		userRepository: userRepository,
	} 
}
