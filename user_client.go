package main

import "fmt"

type UserResponse struct {
	User struct {
		Username string `json:"username"`
	} `json:"user"`
}

func getCurrentUsername() (string, error) {
	url := fmt.Sprintf("%s/current_user.json", apiBaseURL)
	var response UserResponse
	if err := MakeGETRequest(url, &response); err != nil {
		return "", err
	}
	return response.User.Username, nil
}
