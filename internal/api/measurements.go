package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/tehnerd/vape/internal/models"
)

func (c *Client) CreateMeasurement(req *models.MeasurementRequest) (*models.MeasurementCreateResponse, error) {
	body, err := c.Post("/measurements/", req)
	if err != nil {
		return nil, err
	}

	var resp models.MeasurementCreateResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp, nil
}

func (c *Client) GetMeasurement(id int) (*models.Measurement, error) {
	body, err := c.Get(fmt.Sprintf("/measurements/%d/", id), nil)
	if err != nil {
		return nil, err
	}

	var measurement models.Measurement
	if err := json.Unmarshal(body, &measurement); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &measurement, nil
}

type MeasurementListParams struct {
	Status   string
	Type     string
	Target   string
	IsOneoff *bool
	Mine     bool
	Limit    int
	Offset   int
}

func (c *Client) ListMeasurements(params *MeasurementListParams) (*models.MeasurementListResponse, error) {
	query := url.Values{}

	if params != nil {
		if params.Status != "" {
			query.Set("status", params.Status)
		}
		if params.Type != "" {
			query.Set("type", params.Type)
		}
		if params.Target != "" {
			query.Set("target", params.Target)
		}
		if params.IsOneoff != nil {
			query.Set("is_oneoff", strconv.FormatBool(*params.IsOneoff))
		}
		if params.Mine {
			query.Set("mine", "true")
		}
		if params.Limit > 0 {
			query.Set("page_size", strconv.Itoa(params.Limit))
		}
		if params.Offset > 0 {
			query.Set("offset", strconv.Itoa(params.Offset))
		}
	}

	body, err := c.Get("/measurements/", query)
	if err != nil {
		return nil, err
	}

	var resp models.MeasurementListResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp, nil
}

func (c *Client) GetMeasurementResults(id int, start, stop int64, probeIDs []int, limit int) ([]models.MeasurementResult, error) {
	query := url.Values{}

	if start > 0 {
		query.Set("start", strconv.FormatInt(start, 10))
	}
	if stop > 0 {
		query.Set("stop", strconv.FormatInt(stop, 10))
	}
	if len(probeIDs) > 0 {
		for _, pid := range probeIDs {
			query.Add("probe_ids", strconv.Itoa(pid))
		}
	}
	if limit > 0 {
		query.Set("limit", strconv.Itoa(limit))
	}

	body, err := c.Get(fmt.Sprintf("/measurements/%d/results/", id), query)
	if err != nil {
		return nil, err
	}

	var results []models.MeasurementResult
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return results, nil
}

func (c *Client) StopMeasurement(id int) error {
	_, err := c.Delete(fmt.Sprintf("/measurements/%d/", id))
	return err
}
