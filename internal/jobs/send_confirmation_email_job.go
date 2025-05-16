package jobs

import (
	"context"
	"database/sql"
	"log"
	"math"
	"time"

	"github.com/kievzenit/genesis-case/internal/database"
	"github.com/kievzenit/genesis-case/internal/database/repositories"
	"github.com/kievzenit/genesis-case/internal/services"
)

type SendConfirmationEmailJob struct {
	emailService services.EmailService
	txManager    *database.TransactionManger
}

func NewSendConfirmationEmailJob(
	emailService services.EmailService,
	txManager *database.TransactionManger,
) *SendConfirmationEmailJob {
	return &SendConfirmationEmailJob{
		emailService: emailService,
		txManager:    txManager,
	}
}

func (job *SendConfirmationEmailJob) Run() {
	err := job.txManager.ExecuteTx(func(tx *sql.Tx) error {
		confirmationEmailsRepository := repositories.NewConfirmationEmailsRepository(tx)
		subscriptionsRepository := repositories.NewSubscriptionRepository(tx)

		emails, err := confirmationEmailsRepository.GetConfirmationEmailsToSend(context.Background())
		if err != nil {
			return err
		}

		for _, email := range emails {
			subscription, err := subscriptionsRepository.GetSubscriptionByTokenContext(
				context.Background(), 
				email.Token,
			)
			if err != nil {
				return err
			}
						
			err = job.emailService.SendConfirmationEmail(
				email.ToAddress,
				subscription.City,
				subscription.Frequency,
				email.Token,
			)
			if err != nil {
				email.Attempts++

				delay := time.Duration(math.Pow(2, float64(email.Attempts))) * time.Minute
				email.NextTryAfter = time.Now().UTC().Add(delay)

				err = confirmationEmailsRepository.UpdateConfirmationEmail(context.Background(), email)
				if err != nil {
					return err
				}

				continue
			}

			email.Attempts++
			email.Completed = true
			err = confirmationEmailsRepository.UpdateConfirmationEmail(context.Background(), email)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		log.Printf("send confirmation email job finished with error: %v", err)
	}
}
