package repository

import (
	"MedicineBuddy/model"
	"context"
)

type UserRepository interface {
	ListOfUsers(ctx context.Context) ([]model.User, error)
}
