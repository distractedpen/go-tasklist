package util

type ErrDoesNotExist struct{}
func (e ErrDoesNotExist) Error() string {
	return "Entity does not Exist"
}
func (e ErrDoesNotExist) Is(other error) bool {
	return other == e
}

type ErrExists struct {}
func (e ErrExists) Error() string {
	return "Entity Exists"
}
func (e ErrExists) Is(other error) bool {
	return other == e
}

type ErrDB struct {}
func (e ErrDB) Error() string {
	return "Database Error"
}
func (e ErrDB) Is(other error) bool {
	return other == e
}

type ErrRequestInvalid struct {}
func (e ErrRequestInvalid) Error() string {
	return "Invalid Request"
}
func (e ErrRequestInvalid) Is(other error) bool {
	return other == e
}

type ErrAuthInvalid struct {}
func (e ErrAuthInvalid) Error() string {
	return "Invalid Username or Password"
}
func (e ErrAuthInvalid) Is(other error) bool {
	return other == e
}
