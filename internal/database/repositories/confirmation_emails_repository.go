package repositories

import (
	"context"
	"time"

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
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE pending_confirmation_emails
		SET completed = $1, attempts = $2, next_try_after = $3
		WHERE id = $4`,
		confirmationEmail.Completed,
		confirmationEmail.Attempts,
		confirmationEmail.NextTryAfter,
		confirmationEmail.Id,
	)
	return err
}

func (r *confirmationEmailsRepository) GetConfirmationEmailsToSend(
	ctx context.Context,
) ([]models.ConfirmationEmail, error) {
	nowUtc := time.Now().UTC()
	confirmationEmailsRows, err := r.db.QueryContext(
		ctx,
		`SELECT id, to_address, token, attempts, next_try_after
		FROM pending_confirmation_emails
		WHERE completed = false AND next_try_after <= $1 AND attempts < 3
		FOR UPDATE SKIP LOCKED`,
		nowUtc,
	)
	if err != nil {
		return nil, err
	}
	defer confirmationEmailsRows.Close()

	var confirmationEmails []models.ConfirmationEmail
	for confirmationEmailsRows.Next() {
		var confirmationEmail models.ConfirmationEmail
		err := confirmationEmailsRows.Scan(
			&confirmationEmail.Id,
			&confirmationEmail.ToAddress,
			&confirmationEmail.Token,
			&confirmationEmail.Attempts,
			&confirmationEmail.NextTryAfter,
		)
		if err != nil {
			return nil, err
		}
		confirmationEmails = append(confirmationEmails, confirmationEmail)
	}

	return confirmationEmails, nil
}

func (r *confirmationEmailsRepository) StoreConfirmationEmail(
	ctx context.Context,
	confirmationEmail models.ConfirmationEmail,
) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO pending_confirmation_emails (to_address, token, next_try_after)
		VALUES ($1, $2, $3)`,
		confirmationEmail.ToAddress,
		confirmationEmail.Token,
		confirmationEmail.NextTryAfter,
	)
	return err
}
