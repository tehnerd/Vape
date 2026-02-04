package api

import (
	"encoding/json"
	"fmt"

	"github.com/tehnerd/vape/internal/models"
)

func (c *Client) GetCredits() (*models.Credits, error) {
	body, err := c.Get("/credits/", nil)
	if err != nil {
		return nil, err
	}

	var credits models.Credits
	if err := json.Unmarshal(body, &credits); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &credits, nil
}
