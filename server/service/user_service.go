package service

import (
	"MedicineBuddy/model"
	"MedicineBuddy/repository"
	"context"
)

type UserService struct {
	userRepository repository.UserRepository
}

func NewUserService(us repository.UserRepository) *UserService {
	return &UserService{userRepository: us}
}

func (us *UserService) ListOfUsers(ctx context.Context) ([]model.User, error) {
	return us.userRepository.ListOfUsers(ctx)
}
