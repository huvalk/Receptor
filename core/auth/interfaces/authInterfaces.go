package authInterfaces

import "GOSecretProject/core/model/base"

type AuthRepository interface {
	Register(user base.User) (err error)
	Login(user base.User) (userID int, session string, statusCode int)
	Logout(session string) (err error)
}
