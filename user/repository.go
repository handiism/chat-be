package user

import "github.com/handiism/chat/entity"

type Repo interface {
	Store(name, password string) (entity.User, error)
	FetchByUsernameAndPassword(name, password string) (entity.User, error)
	FetchById(id string) (entity.User, error)
}
