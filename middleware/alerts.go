package middleware

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

func (c *Client) GetAlerts(ctx context.Context, ruleID int, params *GetAlertsParams) (*AlertsResponse, error) {
	path := fmt.Sprintf("/rules/%d/alerts", ruleID)
	
	if params != nil {
		query := url.Values{}
		if params.Page > 0 {
			query.Set("page", strconv.Itoa(params.Page))
		}
		if params.OrderBy != "" {
			query.Set("order_by", params.OrderBy)
		}
		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	var result AlertsResponse
	if err := c.doRequest(ctx, "GET", path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) CreateAlert(ctx context.Context, ruleID int, alert *NewAlert) (*Alert, error) {
	path := fmt.Sprintf("/rules/%d/alerts", ruleID)
	var result Alert
	if err := c.doRequest(ctx, "POST", path, alert, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetAlertStats(ctx context.Context, ruleID int) (*StatsResponse, error) {
	path := fmt.Sprintf("/rules/%d/alerts/stats", ruleID)
	var result StatsResponse
	if err := c.doRequest(ctx, "GET", path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

type GetAlertsParams struct {
	Page    int
	OrderBy string
}

