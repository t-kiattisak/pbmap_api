package pokeapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"pbmap_api/src/internal/adapter/models"
	"pbmap_api/src/internal/usecase"
	"time"
)

type PokeAPIClient struct {
	client  *http.Client
	baseURL string
}

func NewPokeAPIClient() *PokeAPIClient {
	return &PokeAPIClient{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: "https://pokeapi.co/api/v2",
	}
}

func (c *PokeAPIClient) FetchDitto(ctx context.Context) (*models.PokemonDittoReponse, error) {
	url := fmt.Sprintf("%s/pokemon/ditto", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var data models.PokemonDittoReponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &data, nil
}

func (c *PokeAPIClient) FetchExternalData(ctx context.Context) (*usecase.ExternalDataStub, error) {
	resp, err := c.FetchDitto(ctx)
	if err != nil {
		return nil, err
	}
	return &usecase.ExternalDataStub{Value: resp.BaseExperience}, nil
}
