package utils

import (
	"fmt"
	"net/url"
	"strings"
)

func EnsureFullUrl(newUrl, fetchedUrl string) (string, error) {
	newUrl = strings.TrimSpace(newUrl)
	if newUrl == "" {
		return "", nil
	}

	// Handle relative URLs starting with "./"
	if strings.HasPrefix(newUrl, "./") {
		newUrl = strings.TrimPrefix(newUrl, "./")
		baseURL, err := url.Parse(fetchedUrl)
		if err != nil {
			return "", fmt.Errorf("error parsing URL: %w", err)
		}

		// Remove the last part of the path from the fetched URL
		lastSlashIndex := strings.LastIndex(baseURL.Path, "/")
		if lastSlashIndex != -1 {
			baseURL.Path = baseURL.Path[:lastSlashIndex]
		}

		// Ensure there is a slash between the modified path and the new URL path
		if !strings.HasSuffix(baseURL.Path, "/") {
			baseURL.Path += "/"
		}

		// Construct the new full URL
		newFullUrl := baseURL.Scheme + "://" + baseURL.Host + baseURL.Path + newUrl
		return newFullUrl, nil
	}

	// Handle absolute URLs that do not start with "http://" or "https://"
	if !strings.HasPrefix(newUrl, "http://") && !strings.HasPrefix(newUrl, "https://") {
		baseURL, err := url.Parse(fetchedUrl)
		if err != nil {
			return "", fmt.Errorf("error parsing URL: %w", err)
		}
		// Ensure there is a slash between the host and the new URL path
		if !strings.HasPrefix(newUrl, "/") {
			newUrl = "/" + newUrl
		}
		return baseURL.Scheme + "://" + baseURL.Host + newUrl, nil
	}

	// Return the new URL if it's already a complete URL
	return newUrl, nil
}
