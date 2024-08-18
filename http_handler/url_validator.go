package http_handler

import (
	"fmt"
	"net/url"
)

func isValidUrl(rawUrl string) (bool, error) {
	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		return false, fmt.Errorf("error parsing URL: %s", err)
	}

	if parsedUrl.Scheme != "http" && parsedUrl.Scheme != "https" {
		return false, fmt.Errorf("scheme must be http or https")
	}

	if parsedUrl.Host == "" {
		return false, fmt.Errorf("host is empty")
	}

	return true, nil
}
