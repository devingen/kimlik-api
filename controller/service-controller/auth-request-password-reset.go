package service_controller

import (
	"context"
	"net/http"
	"time"

	core "github.com/devingen/api-core"
	ulak_dto "github.com/devingen/ulak-api/dto"
	ulak_http_client "github.com/devingen/ulak-api/http-client"
	ulak_model "github.com/devingen/ulak-api/model"
)

type RequestPasswordResetRequest struct {
	Email  string `json:"email" validate:"required,email"`
	Origin string `json:"origin" validate:"required"`
}

func (c ServiceController) RequestPasswordReset(ctx context.Context, req core.Request) (*core.Response, error) {

	base, hasBase := req.PathParameters["base"]
	if !hasBase {
		return nil, core.NewError(http.StatusInternalServerError, "missing-path-param-base")
	}

	var body RequestPasswordResetRequest
	err := req.AssertBody(&body)
	if err != nil {
		return nil, err
	}

	// Find user by email
	user, err := c.DataService.FindUserWithEmail(ctx, base, body.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, core.NewError(http.StatusNotFound, "user-not-found")
	}

	// Check if user has password auth
	auth, err := c.DataService.FindPasswordAuthOfUser(ctx, base, user.ID.Hex())
	if err != nil {
		return nil, err
	}

	if auth == nil {
		return nil, core.NewError(http.StatusBadRequest, "password-auth-not-available")
	}

	// Get integration settings for Ulak configuration
	integrationSettings, err := c.DataService.GetIntegrationSettings(ctx, base)
	if err != nil {
		return nil, err
	}

	if integrationSettings == nil || integrationSettings.Ulak == nil {
		return nil, core.NewError(http.StatusInternalServerError, "ulak-not-configured")
	}

	ulak := integrationSettings.Ulak

	// Generate reset token with "set-password" scope
	resetToken, err := c.TokenService.GenerateAccessToken(
		user.ID.Hex(),
		"",
		[]string{"set-password"},
		time.Now().Add(24*time.Hour).Unix(),
	)
	if err != nil {
		return nil, err
	}

	// Determine reset page URL from request body origin
	resetPageURL := body.Origin + "/reset-password?token=" + resetToken

	// Create Ulak API client
	ulakHeaders := "devingen-product-id=" + *ulak.ProductID + ",api-key=" + *ulak.APIKey
	ulakAPIClient := ulak_http_client.New("https://api.ulak.das.devingen.io", ulakHeaders)

	// Send reset email via Ulak
	_, err = ulakAPIClient.SendEmail(ctx, ulak_dto.SendEmailRequest{
		EmailSenderConfigurationID: *ulak.SenderConfigurationID,
		EmailTemplateIdentifier:    "reset-password",
		Language:                   "en-us",
		Parameters: map[string]interface{}{
			"recipient.name":  *user.FirstName,
			"recipient.email": *user.Email,
			"page.url":        resetPageURL,
		},
		ReplyToEmail: *ulak.SenderEmail,
		Sender: ulak_model.EmailRecipient{
			Name:  *ulak.SenderName,
			Email: *ulak.SenderEmail,
		},
		To: []ulak_model.EmailRecipient{
			{Email: *user.Email},
		},
		Cc:          nil,
		Bcc:         nil,
		Attachments: nil,
	})
	if err != nil {
		return nil, err
	}

	return &core.Response{
		StatusCode: http.StatusOK,
		Body:       map[string]string{"message": "Password reset email sent"},
	}, nil
}
