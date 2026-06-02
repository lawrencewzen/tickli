package api

import (
	"fmt"
	"github.com/sho0pi/tickli/internal/config"
)

func GetClient() (*Client, error) {
	td, err := config.LoadTokenData()
	if err != nil {
		return nil, fmt.Errorf("failed to load token: %w", err)
	}
	if td == nil || td.AccessToken == "" {
		return nil, fmt.Errorf("no token found, please run 'tickli init' first")
	}

	return NewClient(td.AccessToken), nil
}
