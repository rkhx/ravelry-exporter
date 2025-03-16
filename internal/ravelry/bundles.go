package ravelry

import (
	"context"
	"fmt"
	"github.com/rkhx/ravelry-exporter/internal/httpclient"
	"github.com/rkhx/ravelry-exporter/internal/models"
)

type BundleClient struct {
	HTTPClient *httpclient.CustomHTTPClient
}

func (b *BundleClient) GetUserBundles(ctx context.Context, username string) ([]models.Bundle, error) {
	url := fmt.Sprintf("/people/%s/bundles/list.json", username)
	var response models.BundlesResponse
	if err := b.HTTPClient.MakeGETRequest(ctx, url, &response); err != nil {
		return nil, err
	}

	return response.Bundles, nil
}

func (b *BundleClient) GetBundleContent(ctx context.Context, username string, bundleID int) (*models.Bundle, error) {
	url := fmt.Sprintf("/people/%s/bundles/%d.json", username, bundleID)
	var response models.BundleContentResponse
	if err := b.HTTPClient.MakeGETRequest(ctx, url, &response); err != nil {
		return nil, err
	}

	return &response.Bundle, nil
}

func (b *BundleClient) GetBundleItem(ctx context.Context, itemID int) (*models.BundledItemResponse, error) {
	url := fmt.Sprintf("/bundled_items/%d.json", itemID)
	var response models.BundledItemResponse
	if err := b.HTTPClient.MakeGETRequest(ctx, url, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
