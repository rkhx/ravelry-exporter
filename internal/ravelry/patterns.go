package ravelry

import (
	"context"
	"fmt"
	"github.com/rkhx/ravelry-exporter/internal/httpclient"
	"github.com/rkhx/ravelry-exporter/internal/models"
)

type PatternsClient struct {
	HTTPClient *httpclient.CustomHTTPClient
}

func (p *PatternsClient) GetPattern(ctx context.Context, patternID int) (*models.Pattern, error) {
	url := fmt.Sprintf("/patterns/%d.json", patternID)
	var response models.PatternResponse
	if err := p.HTTPClient.MakeGETRequest(ctx, url, &response); err != nil {
		return nil, err
	}

	return &response.Pattern, nil
}

func (p *PatternsClient) GetPatternAttr(ctx context.Context) (interface{}, error) {
	var response interface{}
	if err := p.HTTPClient.MakeGETRequest(ctx, "/yarn_weights.json", &response); err != nil {
		return nil, err
	}

	return response, nil
}
