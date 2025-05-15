package repositories

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/kievzenit/genesis-case/internal/database"
	"github.com/kievzenit/genesis-case/internal/models"
)

type SubscriptionRepository interface {
	IsUserSubscribedContext(ctx context.Context, email string, city string) (bool, error)
	SubscribeContext(
		ctx context.Context,
		email string,
		token uuid.UUID,
		city string,
		frequency models.Frequency,
	) error
	ConfirmSubscriptionContext(ctx context.Context, token uuid.UUID) error
	UnsubscribeContext(ctx context.Context, token uuid.UUID) error
}

func NewSubscriptionRepository(db database.Database) SubscriptionRepository {
	return &subscriptionRepository{
		db: db,
	}
}

type subscriptionRepository struct {
	db database.Database
}

var ErrConfirmationTokenNotFound = errors.New("confirmation token not found")

func (r *subscriptionRepository) ConfirmSubscriptionContext(ctx context.Context, token uuid.UUID) error {
	panic("unimplemented")
}

func (r *subscriptionRepository) IsUserSubscribedContext(ctx context.Context, email string, city string) (bool, error) {
	panic("unimplemented")
}

func (r *subscriptionRepository) SubscribeContext(
	ctx context.Context,
	email string,
	token uuid.UUID,
	city string,
	frequency models.Frequency,
) error {
	panic("unimplemented")
}

func (r *subscriptionRepository) UnsubscribeContext(ctx context.Context, token uuid.UUID) error {
	panic("unimplemented")
}
