package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/tehnerd/vape/internal/models"
)

func (c *Client) GetProbe(id int) (*models.Probe, error) {
	body, err := c.Get(fmt.Sprintf("/probes/%d/", id), nil)
	if err != nil {
		return nil, err
	}

	var probe models.Probe
	if err := json.Unmarshal(body, &probe); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &probe, nil
}

func (c *Client) ListProbes(params *models.ProbeListParams) (*models.ProbeListResponse, error) {
	query := url.Values{}

	if params != nil {
		if params.CountryCode != "" {
			query.Set("country_code", params.CountryCode)
		}
		if params.ASN > 0 {
			query.Set("asn", strconv.Itoa(params.ASN))
		}
		if params.Status != "" {
			query.Set("status", string(params.Status))
		}
		if params.IsAnchor != nil {
			query.Set("is_anchor", strconv.FormatBool(*params.IsAnchor))
		}
		if params.IsPublic != nil {
			query.Set("is_public", strconv.FormatBool(*params.IsPublic))
		}
		if params.Limit > 0 {
			query.Set("page_size", strconv.Itoa(params.Limit))
		}
		if params.Offset > 0 {
			query.Set("offset", strconv.Itoa(params.Offset))
		}
	}

	body, err := c.Get("/probes/", query)
	if err != nil {
		return nil, err
	}

	var resp models.ProbeListResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp, nil
}
