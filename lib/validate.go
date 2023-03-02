package lib

import "net/url"

func ValidateURL(inputURL string) bool {
	// Parse the URL
	_, err := url.ParseRequestURI(inputURL)

	// Check for errors and return result
	if err != nil || inputURL == "" {
		return false
	}
	return true
}
