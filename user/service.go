package user

import (
	"fmt"

	"github.com/handiism/chat/entity"
)

type Service interface {
	Login(username, password string) (entity.User, error)
	Register(username, password string) (entity.User, error)
	FetchById(id string) (entity.User, error)
}

type service struct {
	repo Repo
}

func NewService(repo Repo) Service {
	return &service{
		repo: repo,
	}
}

// FetchById implements Service
func (s *service) FetchById(id string) (entity.User, error) {
	fetchedUser, err := s.repo.FetchById(id)
	if err != nil {
		return entity.User{}, err
	}

	return fetchedUser, nil
}

// Login implements Service
func (s *service) Login(username string, password string) (entity.User, error) {
	userLogin, err := s.repo.FetchByUsernameAndPassword(username, password)
	if err != nil {
		return entity.User{}, err
	}

	return userLogin, nil
}

// Register implements Service
func (s *service) Register(username string, password string) (entity.User, error) {
	if len(password) < 4 {
		return entity.User{}, fmt.Errorf("password length must be >= 4")
	}

	userRegister, err := s.repo.Store(username, password)
	if err != nil {
		return entity.User{}, err
	}
	return userRegister, nil
}
