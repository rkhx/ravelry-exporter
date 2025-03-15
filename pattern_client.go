package main

import "fmt"

type Pattern struct {
	Name             string `json:"name"`
	GaugeDescription string `json:"gauge_description"`
}

type PatternResponse struct {
	Pattern Pattern `json:"pattern"`
}

func GetPattern(paternId int) (*PatternResponse, error) {
	url := fmt.Sprintf("%s/patterns/%d.json", apiBaseURL, paternId)
	var response PatternResponse
	if err := MakeGETRequest(url, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
