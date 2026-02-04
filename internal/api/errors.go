package api

import (
	"encoding/json"
	"fmt"
)

type APIError struct {
	StatusCode int
	Message    string
	Detail     string
	Errors     []ErrorDetail
}

type ErrorDetail struct {
	Source string `json:"source"`
	Detail string `json:"detail"`
}

type errorResponse struct {
	Error struct {
		Status  int           `json:"status"`
		Code    int           `json:"code"`
		Title   string        `json:"title"`
		Detail  string        `json:"detail"`
		Errors  []ErrorDetail `json:"errors"`
	} `json:"error"`
}

func (e *APIError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("API error %d: %s - %s", e.StatusCode, e.Message, e.Detail)
	}
	if len(e.Errors) > 0 {
		msg := fmt.Sprintf("API error %d: %s", e.StatusCode, e.Message)
		for _, err := range e.Errors {
			msg += fmt.Sprintf("\n  - %s: %s", err.Source, err.Detail)
		}
		return msg
	}
	return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Message)
}

func ParseAPIError(statusCode int, body []byte) error {
	var errResp errorResponse
	if err := json.Unmarshal(body, &errResp); err != nil {
		return &APIError{
			StatusCode: statusCode,
			Message:    string(body),
		}
	}

	return &APIError{
		StatusCode: statusCode,
		Message:    errResp.Error.Title,
		Detail:     errResp.Error.Detail,
		Errors:     errResp.Error.Errors,
	}
}

func IsNotFound(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == 404
	}
	return false
}

func IsUnauthorized(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == 401
	}
	return false
}

func IsForbidden(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == 403
	}
	return false
}
