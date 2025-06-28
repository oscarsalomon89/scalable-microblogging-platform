package user

import "context"

type (
	Repository interface {
		CreateUser(ctx context.Context, user *User) error
	}

	userUseCase struct {
		repo Repository
	}
)

func NewUserUseCase(repo Repository) *userUseCase {
	return &userUseCase{repo: repo}
}

func (uc *userUseCase) CreateUser(ctx context.Context, user *User) error {
	//TODO: validar username o email repetido
	return uc.repo.CreateUser(ctx, user)
}
