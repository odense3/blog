package service

import (
	"blog/internal/adapter/repository"
	"blog/internal/core/domain/entity"
	"blog/lib/conv"
	"context"

	"github.com/gofiber/fiber/v2/log"
)

type UserService interface {
	UpdatePassword(ctx context.Context, newPass string, id int64) error
	GetUserById(ctx context.Context, id int64) (*entity.UserEntity, error)
}

type userService struct {
	userRepository repository.UserRepository
}

// GetUserById implements UserService.
func (u *userService) GetUserById(ctx context.Context, id int64) (*entity.UserEntity, error) {
	result, err := u.userRepository.GetUserById(ctx, id)
	if err != nil {
		code := "[SERVICE] GetUserById - 1"
		log.Errorw(code, err)
		return nil, err
	}

	return result, nil
}

// UpdatePassword implements UserService.
func (u *userService) UpdatePassword(ctx context.Context, newPass string, id int64) error {
	password, err := conv.HashPassword(newPass)
	if err != nil {
		code := "[SERVICE] UpdatePassword - 1"
		log.Errorw(code, err)
		return err
	}

	err = u.userRepository.UpdatePassword(ctx, password, id)
	if err != nil {
		code := "[SERVICE] UpdatePassword - 2"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}
