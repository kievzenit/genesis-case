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

type subscriptionData struct {
	Email     string `json:"email" form:"email"`
	City      string `json:"city" form:"city"`
	Frequency string `json:"frequency" form:"frequency"`
}

func SubscribeForWeatherHandler(
	weatherService services.WeatherService,
	emailService services.EmailService,
	database database.Database,
	txManager *database.TransactionManger,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		var data subscriptionData

		// In swagger specification was defined that this endpoint should accept both
		// application/json and application/x-www-form-urlencoded content types.
		// Maybe that was an error, but it was said, that we need 100% be complaint with API specification.
		// So I implemented both content types.
		contentType := c.Request.Header.Get("Content-Type")
		if contentType == "application/json" {
			err := c.BindJSON(&data)
			if err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
		} else if contentType == "application/x-www-form-urlencoded" {
			data.Email = c.PostForm("email")
			data.City = c.PostForm("city")
			data.Frequency = c.PostForm("frequency")
		} else {
			c.AbortWithStatus(http.StatusUnsupportedMediaType)
			return
		}

		if data.Email == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		if data.City == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		if data.Frequency == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		frequency := models.Frequency(data.Frequency)
		if !frequency.IsValid() {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		subscriptionRepository := repositories.NewSubscriptionRepository(database)
		exists, err := subscriptionRepository.IsUserSubscribedContext(ctx, data.Email, data.City)
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
			err = subscriptionRepository.SubscribeContext(ctx, data.Email, token, data.City, frequency)
			if err != nil {
				return err
			}

			return confirmationEmailRepository.StoreConfirmationEmail(
				ctx,
				models.ConfirmationEmail{
					ToAddress:    data.Email,
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

func ConfirmSubscriptionHandler(database database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		tokenParam := c.Param("token")
		if tokenParam == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		token, err := uuid.Parse(tokenParam)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		subscriptionRepository := repositories.NewSubscriptionRepository(database)
		err = subscriptionRepository.ConfirmSubscriptionContext(ctx, token)
		if err != nil {
			if err == repositories.ErrConfirmationTokenNotFound {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
}

func UnsubscribeHandler(database database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		tokenParam := c.Param("token")
		if tokenParam == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		token, err := uuid.Parse(tokenParam)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		subscriptionRepository := repositories.NewSubscriptionRepository(database)
		err = subscriptionRepository.UnsubscribeContext(ctx, token)
		if err != nil {
			if err == repositories.ErrConfirmationTokenNotFound {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
}
