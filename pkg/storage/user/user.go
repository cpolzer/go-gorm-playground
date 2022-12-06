package user

import "context"

type Repository interface {
	Has(ctx context.Context, userId uint) (bool, error)
	Get(ctx context.Context, userId uint) (*UserDto, error)
	Create(ctx context.Context, user *UserDto) (*UserDto, error)
	Update(ctx context.Context, project *UserDto) (*UserDto, error)
	Delete(ctx context.Context, userId uint) error
}
