package user

import (
	"context"
	"log"

	"github.com/michael870917/gocourse_domain/domain"
)

type (
	// Structs
	Filters struct {
		FirstName string
		LastName  string
	}

	ServiceInterface interface {
		Create(ctx context.Context, firstName, lastName, email, phone string) (*domain.User, error)
		GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.User, error)
		GetById(ctx context.Context, id string) (*domain.User, error)
		Delete(ctx context.Context, id string) error
		Update(ctx context.Context, id string, firstName, lastName, email, phone *string) error
		Count(ctx context.Context, filters Filters) (int64, error)
	}

	serviceLogic struct {
		logger     *log.Logger
		repository Repository
	}
)

func NewService(logger *log.Logger, repository Repository) ServiceInterface {
	return &serviceLogic{repository: repository, logger: logger}
}

func (sl serviceLogic) Create(ctx context.Context, firstName, lastName, email, phone string) (*domain.User, error) {
	user := domain.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
	}

	userData, error := sl.repository.Create(ctx, &user)
	if error != nil {
		return nil, error
	}

	return userData, nil
}

func (sl serviceLogic) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.User, error) {
	users, error := sl.repository.GetAll(ctx, filters, offset, limit)
	if error != nil {
		return nil, error
	}
	return users, nil
}

func (sl serviceLogic) GetById(ctx context.Context, id string) (*domain.User, error) {
	user, error := sl.repository.GetByID(ctx, id)
	if error != nil {
		return nil, error
	}
	return user, nil
}

func (sl serviceLogic) Delete(ctx context.Context, id string) error {
	return sl.repository.Delete(ctx, id)
}

func (sl serviceLogic) Update(ctx context.Context, id string, firstName, lastName, email, phone *string) error {
	return sl.repository.Update(ctx, id, firstName, lastName, email, phone)
}

func (sl serviceLogic) Count(ctx context.Context, filters Filters) (int64, error) {
	return sl.repository.Count(ctx, filters)
}
