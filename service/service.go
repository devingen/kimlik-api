package service

import "github.com/devingen/kimlik-api/model"

type IKimlikService interface {
	CreateAuth(base, password string, user *model.User) (*model.Auth, error)
	CreateSession(base, client, userAgent, ip string, user *model.User) (*model.Session, error)
	CreateUser(base, firstName, lastName, email string) (*model.User, error)
	FindAuthOfUser(base string, user *model.User, authType model.AuthType) (*model.Auth, error)
	FindUserUserWithEmail(base, email string) (*model.User, error)
	FindUserUserWithId(base, id string) (*model.User, error)
}
