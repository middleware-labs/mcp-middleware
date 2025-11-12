package middleware

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

func (c *Client) GetWidgets(ctx context.Context, params *GetWidgetsParams) (*Widget, error) {
	path := "/builder/widget"
	
	if params != nil {
		query := url.Values{}
		if params.ReportID > 0 {
			query.Set("report_id", strconv.Itoa(params.ReportID))
		}
		if params.DisplayScope != "" {
			query.Set("display_scope", params.DisplayScope)
		}
		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	var result Widget
	if err := c.doRequest(ctx, "GET", path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) CreateWidget(ctx context.Context, widget *CustomWidget) (*Widget, error) {
	var result Widget
	if err := c.doRequest(ctx, "POST", "/builder/widget", widget, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) DeleteWidget(ctx context.Context, builderID int) error {
	path := fmt.Sprintf("/builder/widget/%d", builderID)
	return c.doRequest(ctx, "DELETE", path, nil, nil)
}

func (c *Client) GetWidgetData(ctx context.Context, widget *CustomWidget) (*BuilderDataResponse, error) {
	var result BuilderDataResponse
	if err := c.doRequest(ctx, "POST", "/builder/widget/data", widget, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetMultiWidgetData(ctx context.Context, widgets []CustomWidget) ([]BuilderDataResponse, error) {
	var result []BuilderDataResponse
	if err := c.doRequest(ctx, "POST", "/builder/widget/multi-data", widgets, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) UpdateWidgetLayouts(ctx context.Context, req *LayoutRequest) error {
	return c.doRequest(ctx, "PUT", "/builder/widget/scope/layouts", req, nil)
}

type GetWidgetsParams struct {
	ReportID     int
	DisplayScope string
}

