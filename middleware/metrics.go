package middleware

import (
	"context"
)

func (c *Client) GetMetrics(ctx context.Context, req *MetricsV2Request) (*MetricsV2Response, error) {
	var result MetricsV2Response
	if err := c.doRequest(ctx, "POST", "/builder/metrics-v2", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetResources(ctx context.Context) ([]string, error) {
	var result []string
	if err := c.doRequest(ctx, "GET", "/builder/resources", nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) Query(ctx context.Context, req *QueryRequest) (*QueryResponse, error) {
	var result QueryResponse
	if err := c.doRequest(ctx, "POST", "/query", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
