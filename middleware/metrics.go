package middleware

import (
	"context"
)

// GetMetrics retrieves metrics, filters, or groupby tags
func (c *Client) GetMetrics(ctx context.Context, req *MetricsV2Request) (*MetricsV2Response, error) {
	var result MetricsV2Response
	if err := c.doRequest(ctx, "POST", "/builder/metrics-v2", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetResources retrieves a list of available resources
func (c *Client) GetResources(ctx context.Context) ([]string, error) {
	var result []string
	if err := c.doRequest(ctx, "GET", "/builder/resources", nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

