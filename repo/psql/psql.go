package psql

import (
	"context"
	"log"
	"time"

	"github.com/gofrs/uuid"
	"github.com/handiism/chat/entity"
	"github.com/handiism/chat/user"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) user.Repo {
	return &repo{
		pool: pool,
	}
}

func (r *repo) Store(username, password string) (entity.User, error) {
	id := uuid.UUID{}
	createdAt := time.Time{}
	newUser := entity.User{}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Println(err.Error())
		return entity.User{}, err
	}

	if err := r.pool.QueryRow(context.Background(),
		"INSERT INTO users(username,password) VALUES ($1,$2) RETURNING id, created_at", username, hashedPassword,
	).Scan(&id, &createdAt); err != nil {
		log.Println(err.Error())
		return entity.User{}, err
	}

	newUser.Id = id
	newUser.Username = username
	newUser.Password = password
	newUser.CreatedAt = createdAt

	return newUser, nil
}

func (r *repo) FetchById(id string) (entity.User, error) {
	fetchedUser := entity.User{}

	if err := r.pool.QueryRow(context.Background(),
		"SELECT id, username, password, created_at FROM users WHERE id = $1", id,
	).Scan(&fetchedUser.Id, &fetchedUser.Username, &fetchedUser.Password, &fetchedUser.CreatedAt); err != nil {
		log.Println(err.Error())
		return entity.User{}, nil
	}

	return fetchedUser, nil
}

func (r *repo) FetchByUsernameAndPassword(username, password string) (entity.User, error) {
	fetchedUser := entity.User{}
	hashedPassword := ""

	if err := r.pool.QueryRow(context.Background(),
		"SELECT id, username, password, created_at FROM users WHERE username = $1", username,
	).Scan(&fetchedUser.Id, &fetchedUser.Username, &hashedPassword, &fetchedUser.CreatedAt); err != nil {
		log.Println(err.Error())
		return entity.User{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		log.Println(err.Error())
		return entity.User{}, err
	}

	fetchedUser.Password = hashedPassword

	return fetchedUser, nil
}
