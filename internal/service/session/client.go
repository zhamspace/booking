package session

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	baseURL string
	http    *http.Client
}

func New(baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type ListParticipatedReq struct {
	UserID           string
	Page             int64
	PageSize         int64
	Sort             string
	From             *time.Time
	To               *time.Time
	IncludeCancelled bool
}

type ParticipatedSession struct {
	ID                string     `json:"id"`
	SessionID         string     `json:"session_id"`
	VenueID           string     `json:"venue_id"`
	ResourceID        string     `json:"resource_id"`
	UserID            string     `json:"user_id"`
	Name              string     `json:"name"`
	Status            string     `json:"status"`
	ParticipantStatus string     `json:"participant_status"`
	PaymentStatus     string     `json:"payment_status"`
	PriceTotal        int64      `json:"price_total"`
	Currency          string     `json:"currency"`
	StartAt           *time.Time `json:"start_at"`
	EndAt             *time.Time `json:"end_at"`
	Timezone          string     `json:"timezone"`
	CreatedAt         *time.Time `json:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at"`
	CancelledAt       *time.Time `json:"cancelled_at"`
	PaymentIntentID   *string    `json:"payment_intent_id"`
}

type ListParticipatedRep struct {
	Results        []ParticipatedSession `json:"results"`
	PaginationInfo struct {
		Page       int64 `json:"page"`
		PageSize   int64 `json:"page_size"`
		TotalCount int64 `json:"total_count"`
	} `json:"pagination_info"`
}

func (c *Client) ListParticipated(ctx context.Context, req ListParticipatedReq) (*ListParticipatedRep, error) {
	if strings.TrimSpace(req.UserID) == "" {
		return nil, fmt.Errorf("user_id is required")
	}
	q := url.Values{}
	q.Set("list_params.page", fmt.Sprint(req.Page))
	q.Set("list_params.page_size", fmt.Sprint(req.PageSize))
	if req.Sort != "" {
		q.Set("list_params.sort", req.Sort)
	}
	q.Set("include_cancelled", fmt.Sprint(req.IncludeCancelled))
	if req.From != nil {
		q.Set("from", req.From.UTC().Format(time.RFC3339))
	}
	if req.To != nil {
		q.Set("to", req.To.UTC().Format(time.RFC3339))
	}

	endpoint := c.baseURL + "/internal/users/" + url.PathEscape(req.UserID) + "/sessions?" + q.Encode()
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("session list participated returned %s", resp.Status)
	}

	var out ListParticipatedRep
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}
