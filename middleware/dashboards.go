package middleware

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// GetDashboards retrieves a list of dashboards/reports
func (c *Client) GetDashboards(ctx context.Context, params *GetDashboardsParams) (*ReportListResponse, error) {
	path := "/builder/report"
	
	if params != nil {
		query := url.Values{}
		if params.Limit > 0 {
			query.Set("limit", strconv.Itoa(params.Limit))
		}
		if params.Offset > 0 {
			query.Set("offset", strconv.Itoa(params.Offset))
		}
		if params.Search != "" {
			query.Set("search", params.Search)
		}
		if params.FilterBy != "" {
			query.Set("filterBy", params.FilterBy)
		}
		if params.DisplayScope != "" {
			query.Set("display_scope", params.DisplayScope)
		}
		if params.Sort != "" {
			query.Set("sort", params.Sort)
		}
		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	var result ReportListResponse
	if err := c.doRequest(ctx, "GET", path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetDashboardByKey retrieves a specific dashboard by its key
func (c *Client) GetDashboardByKey(ctx context.Context, reportKey string) (*ReportListResponse, error) {
	path := fmt.Sprintf("/builder/report/%s", reportKey)
	
	var result ReportListResponse
	if err := c.doRequest(ctx, "GET", path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateDashboard creates a new dashboard
func (c *Client) CreateDashboard(ctx context.Context, req *UpsertReportRequest) (*Report, error) {
	var result Report
	if err := c.doRequest(ctx, "POST", "/builder/report", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateDashboard updates an existing dashboard
func (c *Client) UpdateDashboard(ctx context.Context, id int, req *UpsertReportRequest) (*Report, error) {
	path := fmt.Sprintf("/builder/report/%d", id)
	var result Report
	if err := c.doRequest(ctx, "PUT", path, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteDashboard deletes a dashboard
func (c *Client) DeleteDashboard(ctx context.Context, id int) error {
	path := fmt.Sprintf("/builder/report/%d", id)
	return c.doRequest(ctx, "DELETE", path, nil, nil)
}

// CloneDashboard clones an existing dashboard
func (c *Client) CloneDashboard(ctx context.Context, req *UpsertReportRequest) (*Report, error) {
	var result Report
	if err := c.doRequest(ctx, "POST", "/builder/report/clone", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SetDashboardFavorite marks a dashboard as favorite or unfavorite
func (c *Client) SetDashboardFavorite(ctx context.Context, reportID int, favorite bool) error {
	path := fmt.Sprintf("/builder/report/favourite/%d/%t", reportID, favorite)
	return c.doRequest(ctx, "GET", path, nil, nil)
}

// GetDashboardsParams contains parameters for listing dashboards
type GetDashboardsParams struct {
	Limit        int
	Offset       int
	Search       string
	FilterBy     string
	DisplayScope string
	Sort         string
}

