package service

import "github.com/devingen/kimlik-api/model"

type KimlikService interface {
	CreateAuth(base, password string, user *model.User) (*model.Auth, error)
	CreateSession(base, client, userAgent, ip string, user *model.User) (*model.Session, error)
	CreateUser(base, firstName, lastName, email string) (*model.User, error)
	FindUserUserWithEmail(base, email string) ([]*model.User, error)
}
