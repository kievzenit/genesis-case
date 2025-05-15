package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kievzenit/genesis-case/internal/database"
	"github.com/kievzenit/genesis-case/internal/database/repositories"
	"github.com/kievzenit/genesis-case/internal/models"
	"github.com/kievzenit/genesis-case/internal/services"
)

func SubscribeForWeatherHandler(
	weatherService services.WeatherService,
	emailService services.EmailService,
	database database.Database,
	txManager *database.TransactionManger,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		email := c.PostForm("email")
		if email == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		city := c.PostForm("city")
		if city == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		frequencyString := c.PostForm("frequency")
		if frequencyString == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		frequency := models.Frequency(frequencyString)
		if !frequency.IsValid() {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		subscriptionRepository := repositories.NewSubscriptionRepository(database)
		exists, err := subscriptionRepository.IsUserSubscribedContext(ctx, email, city)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if exists {
			c.AbortWithStatus(http.StatusConflict)
			return
		}

		err = txManager.ExecuteTx(func(tx *sql.Tx) error {
			subscriptionRepository = repositories.NewSubscriptionRepository(tx)
			confirmationEmailRepository := repositories.NewConfirmationEmailsRepository(tx)

			token := uuid.New()
			err = subscriptionRepository.SubscribeContext(ctx, email, token, city, frequency)
			if err != nil {
				return err
			}

			return confirmationEmailRepository.StoreConfirmationEmail(
				ctx,
				models.ConfirmationEmail{
					ToAddress:    email,
					Token:        token,
					NextTryAfter: time.Now().UTC(),
				},
			)
		})
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
}
