package jobs

import (
	"context"
	"database/sql"
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
	job.txManager.ExecuteTx(func(tx *sql.Tx) error {
		repository := repositories.NewConfirmationEmailsRepository(tx)

		emails, err := repository.GetConfirmationEmailsToSend(context.Background())
		if err != nil {
			return err
		}

		for _, email := range emails {
			err = job.emailService.SendConfirmationEmail(email.ToAddress, email.Token)
			if err != nil {
				email.Attempts++

				delay := time.Duration(math.Pow(2, float64(email.Attempts))) * time.Minute
				email.NextTryAfter = time.Now().UTC().Add(delay)

				err = repository.UpdateConfirmationEmail(context.Background(), email)
				if err != nil {
					return err
				}

				continue
			}

			email.Attempts++
			email.Completed = true
			err = repository.UpdateConfirmationEmail(context.Background(), email)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
