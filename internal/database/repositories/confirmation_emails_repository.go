package repositories

import (
	"context"

	"github.com/kievzenit/genesis-case/internal/database"
	"github.com/kievzenit/genesis-case/internal/models"
)

type ConfirmationEmailsRepository interface {
	StoreConfirmationEmail(context.Context, models.ConfirmationEmail) error
	GetConfirmationEmailsToSend(context.Context) ([]models.ConfirmationEmail, error)
	UpdateConfirmationEmail(context.Context, models.ConfirmationEmail) error
}

func NewConfirmationEmailsRepository(db database.Database) ConfirmationEmailsRepository {
	return &confirmationEmailsRepository{db: db}
}

type confirmationEmailsRepository struct {
	db database.Database
}

func (r *confirmationEmailsRepository) UpdateConfirmationEmail(
	ctx context.Context,
	confirmationEmail models.ConfirmationEmail,
) error {
	panic("unimplemented")
}

func (r *confirmationEmailsRepository) GetConfirmationEmailsToSend(
	ctx context.Context,
) ([]models.ConfirmationEmail, error) {
	panic("unimplemented")
}

func (r *confirmationEmailsRepository) StoreConfirmationEmail(
	ctx context.Context,
	confirmationEmail models.ConfirmationEmail,
) error {
	panic("unimplemented")
}
