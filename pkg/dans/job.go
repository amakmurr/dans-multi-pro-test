package dans

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
)

type Job struct {
	Company     string    `json:"company"`
	CompanyLogo string    `json:"company_logo"`
	CompanyUrl  string    `json:"company_url"`
	CreatedAt   string    `json:"created_at"`
	Description string    `json:"description"`
	HowToApply  string    `json:"how_to_apply"`
	Id          uuid.UUID `json:"id"`
	Location    string    `json:"location"`
	Title       string    `json:"title"`
	Type        string    `json:"type"`
	Url         string    `json:"url"`
}

type GetJobListParams struct {
	Description *string
	Location    *string
	FullTime    *bool
	Page        *int
}

func (c *Client) GetJobList(ctx context.Context, params *GetJobListParams) ([]*Job, error) {
	u, err := url.Parse(fmt.Sprintf("%s/recruitment/positions.json", c.baseURL.String()))
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	if params != nil {
		q := httpReq.URL.Query()
		if params.Description != nil {
			q.Add("description", *params.Description)
		}
		if params.Location != nil {
			q.Add("location", *params.Location)
		}
		if params.FullTime != nil && *params.FullTime {
			q.Add("full_time", "true")
		}
		if params.Page != nil {
			q.Add("page", fmt.Sprintf("%d", *params.Page))
		}
		httpReq.URL.RawQuery = q.Encode()
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var jobs []*Job
	err = json.NewDecoder(resp.Body).Decode(&jobs)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

func (c *Client) GetJobDetail(ctx context.Context, id uuid.UUID) (*Job, error) {
	u, err := url.Parse(fmt.Sprintf("%s/recruitment/positions/%s", c.baseURL.String(), id.String()))
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	var job Job
	err = json.NewDecoder(resp.Body).Decode(&job)
	if err != nil {
		return nil, err
	}
	return &job, nil
}
