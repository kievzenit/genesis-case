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
	GetSubscriptionByTokenContext(
		ctx context.Context,
		token uuid.UUID,
	) (models.Subscription, error)
	GetConfirmedSubscriptionsByFrequencyContext(
		ctx context.Context,
		frequency models.Frequency,
	) ([]models.Subscription, error)
}

func NewSubscriptionRepository(db database.Database) SubscriptionRepository {
	return &subscriptionRepository{db}
}

type subscriptionRepository struct {
	db database.Database
}

func (r *subscriptionRepository) GetSubscriptionByTokenContext(
	ctx context.Context,
	token uuid.UUID,
) (models.Subscription, error) {
	subscriptionRow := r.db.QueryRowContext(
		ctx,
		`SELECT id, confirmed, email, city, frequency_id
		FROM user_subscriptions
		WHERE token = $1
		LIMIT 1`,
		token,
	)

	var subscription models.Subscription
	var frequencyId int
	err := subscriptionRow.Scan(
		&subscription.Id,
		&subscription.Confirmed,
		&subscription.Email,
		&subscription.City,
		&frequencyId,
	)
	if err != nil {
		return models.Subscription{}, err
	}

	frequencyRow := r.db.QueryRowContext(
		ctx,
		"SELECT name FROM frequencies WHERE id = $1 LIMIT 1",
		frequencyId,
	)
	var frequencyName string
	err = frequencyRow.Scan(&frequencyName)
	if err != nil {
		return models.Subscription{}, err
	}

	subscription.Token = token
	subscription.Frequency = models.Frequency(frequencyName)
	return subscription, nil
}

func (r *subscriptionRepository) GetConfirmedSubscriptionsByFrequencyContext(
	ctx context.Context,
	frequency models.Frequency,
) ([]models.Subscription, error) {
	frequencyRow := r.db.QueryRowContext(
		ctx,
		"SELECT id FROM frequencies WHERE name = $1 LIMIT 1",
		frequency,
	)
	var frequencyId int
	err := frequencyRow.Scan(&frequencyId)
	if err != nil {
		return nil, err
	}

	subscriptionRows, err := r.db.QueryContext(
		ctx,
		`SELECT id, token, confirmed, email, city
		FROM user_subscriptions
		WHERE frequency_id = $1 AND confirmed = true`,
		frequencyId,
	)
	if err != nil {
		return nil, err
	}
	defer subscriptionRows.Close()

	var subscriptions []models.Subscription
	for subscriptionRows.Next() {
		var subscription models.Subscription
		err := subscriptionRows.Scan(
			&subscription.Id,
			&subscription.Token,
			&subscription.Confirmed,
			&subscription.Email,
			&subscription.City,
		)
		if err != nil {
			return nil, err
		}
		subscription.Frequency = frequency
		subscriptions = append(subscriptions, subscription)
	}

	return subscriptions, nil
}

var ErrConfirmationTokenNotFound = errors.New("confirmation token not found")

func (r *subscriptionRepository) ConfirmSubscriptionContext(ctx context.Context, token uuid.UUID) error {
	_, err := r.db.ExecContext(
		ctx,
		"UPDATE user_subscriptions SET confirmed = true WHERE token = $1",
		token,
	)
	return err
}

func (r *subscriptionRepository) IsUserSubscribedContext(ctx context.Context, email string, city string) (bool, error) {
	var exists bool

	err := r.db.QueryRowContext(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM user_subscriptions WHERE email = $1 AND city = $2 LIMIT 1)",
		email,
		city,
	).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *subscriptionRepository) SubscribeContext(
	ctx context.Context,
	email string,
	token uuid.UUID,
	city string,
	frequency models.Frequency,
) error {
	frequencyRow := r.db.QueryRowContext(
		ctx,
		"SELECT id FROM frequencies WHERE name = $1 LIMIT 1",
		frequency,
	)
	var frequencyId int
	err := frequencyRow.Scan(&frequencyId)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(
		ctx,
		"INSERT INTO user_subscriptions (email, token, city, frequency_id) VALUES ($1, $2, $3, $4)",
		email,
		token,
		city,
		frequencyId,
	)
	return err
}

func (r *subscriptionRepository) UnsubscribeContext(ctx context.Context, token uuid.UUID) error {
	_, err := r.db.ExecContext(
		ctx,
		"DELETE FROM user_subscriptions WHERE token = $1",
		token,
	)
	return err
}
