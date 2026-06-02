package utils

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sho0pi/tickli/internal/api"
	"github.com/sho0pi/tickli/internal/config"
)

func LoadClient() api.Client {
	td, err := config.LoadTokenData()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load token. Please run 'tickli init' first")
	}
	if td == nil || td.AccessToken == "" {
		log.Fatal().Msg("No token found. Please run 'tickli init' first")
	}

	// Refresh if token expires within 5 minutes and we have credentials for it
	if td.RefreshToken != "" && td.ClientID != "" && !td.ExpiresAt.IsZero() &&
		time.Now().After(td.ExpiresAt.Add(-5*time.Minute)) {
		log.Info().Msg("Access token expiring, refreshing...")
		refreshed, err := api.RefreshAccessToken(td.ClientID, td.ClientSecret, td.RefreshToken)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to refresh token. Please run 'tickli init' again")
		}
		td.AccessToken = refreshed.AccessToken
		if refreshed.RefreshToken != "" {
			td.RefreshToken = refreshed.RefreshToken
		}
		if refreshed.ExpiresIn > 0 {
			td.ExpiresAt = time.Now().Add(time.Duration(refreshed.ExpiresIn) * time.Second)
		}
		if err := config.SaveTokenData(td); err != nil {
			log.Warn().Err(err).Msg("Failed to persist refreshed token")
		}
	}

	return *api.NewClient(td.AccessToken)
}
