package main

import "fmt"

type BundledItemPattern struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type BundledItemResponse struct {
	BundledItem BundledItem          `json:"bundled_item"`
	Item        []BundledItemPattern `json:"item"`
}

type Bundle struct {
	ID           int           `json:"id"`
	Name         string        `json:"name"`
	BundledItems []BundledItem `json:"bundled_items"`
}

type BundledItem struct {
	ItemName string `json:"item_type"`
	Id       int    `json:"id"`
	ItemId   int    `json:"item_id"`
}

type BundlesResponse struct {
	Bundles []Bundle `json:"bundles"`
}

type BundleContentResponse struct {
	Bundle Bundle `json:"bundle"`
}

func getBundleContent(username string, bundleID int) (*Bundle, error) {
	url := fmt.Sprintf("%s/people/%s/bundles/%d.json", apiBaseURL, username, bundleID)
	var response BundleContentResponse
	if err := MakeGETRequest(url, &response); err != nil {
		return nil, err
	}

	return &response.Bundle, nil
}

func getUserBundles(username string) ([]Bundle, error) {
	url := fmt.Sprintf("%s/people/%s/bundles/list.json", apiBaseURL, username)
	var response BundlesResponse
	if err := MakeGETRequest(url, &response); err != nil {
		return nil, err
	}

	return response.Bundles, nil
}

func getBundleItem(itemId int) (*BundledItemResponse, error) {
	url := fmt.Sprintf("%s/bundled_items/%d.json", apiBaseURL, itemId)
	var response BundledItemResponse
	if err := MakeGETRequest(url, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
