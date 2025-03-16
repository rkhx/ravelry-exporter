package ravelry

import (
	"context"
	"github.com/rkhx/ravelry-exporter/internal/httpclient"
	"github.com/rkhx/ravelry-exporter/internal/models"
)

type UserClient struct {
	HTTPClient *httpclient.CustomHTTPClient
}

// UserResponse Describes response with user info
type UserResponse struct {
	User models.User `json:"user"`
}

func (u *UserClient) GetCurrentUsername(ctx context.Context) (string, error) {
	var response UserResponse
	if err := u.HTTPClient.MakeGETRequest(ctx, "/current_user.json", &response); err != nil {
		return "", err
	}
	return response.User.Username, nil
}
