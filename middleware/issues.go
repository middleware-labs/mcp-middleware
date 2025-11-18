package middleware

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

type GetIncidentsParams struct {
	FromTs int64
	ToTs   int64
	Page   int
	Filter string
	Status string
	Search string
}

func (c *Client) GetIncidents(ctx context.Context, params *GetIncidentsParams) (*IncidentsResponse, error) {
	path := "/ops-ai/incidents"

	if params != nil {
		query := url.Values{}
		if params.FromTs > 0 {
			query.Set("from_ts", strconv.FormatInt(params.FromTs, 10))
		}
		if params.ToTs > 0 {
			query.Set("to_ts", strconv.FormatInt(params.ToTs, 10))
		}
		if params.Page > 0 {
			query.Set("page", strconv.Itoa(params.Page))
		}
		if params.Filter != "" {
			query.Set("filter", params.Filter)
		}
		if params.Status != "" {
			query.Set("status", params.Status)
		}
		if params.Search != "" {
			query.Set("search", params.Search)
		}
		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	var result IncidentsResponse
	if err := c.doRequest(ctx, "GET", path, nil, &result); err != nil {
		return nil, err
	}

	for i := range result.Items {
		if result.Items[i].Fingerprint != "" {
			result.Items[i].IssueURL = fmt.Sprintf("%s/ops-ai?fingerprint=%s", c.baseURL, result.Items[i].Fingerprint)
		}
	}

	return &result, nil
}

type GetIncidentDetailParams struct {
	Fingerprint string
	FromTs      int64
	ToTs        int64
	Filter      string
}

func (c *Client) GetIncidentDetail(ctx context.Context, params *GetIncidentDetailParams) (map[string]any, error) {
	path := "/ops-ai/incident-detail"

	if params != nil {
		query := url.Values{}
		if params.Fingerprint != "" {
			query.Set("fingerprint", params.Fingerprint)
		}
		if params.FromTs > 0 {
			query.Set("from_ts", strconv.FormatInt(params.FromTs, 10))
		}
		if params.ToTs > 0 {
			query.Set("to_ts", strconv.FormatInt(params.ToTs, 10))
		}
		if params.Filter != "" {
			query.Set("filter", params.Filter)
		}
		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	var result map[string]any
	if err := c.doRequest(ctx, "GET", path, nil, &result); err != nil {
		return nil, err
	}

	return result, nil
}
