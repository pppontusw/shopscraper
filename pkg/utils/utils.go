package utils

import (
	"fmt"
	"net/url"
	"strings"
)

func EnsureFullUrl(newUrl, fetchedUrl string, uniqueParameters []string, removeFragment bool) (string, error) {
	newUrl = strings.TrimSpace(newUrl)
	if newUrl == "" {
		return "", nil
	}

	var finalUrl *url.URL
	var err error

	if strings.HasPrefix(newUrl, "./") {
		// Handle relative URLs starting with "./"
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
		finalUrl, err = url.Parse(baseURL.Scheme + "://" + baseURL.Host + baseURL.Path + newUrl)
		if err != nil {
			return "", fmt.Errorf("error constructing final URL: %w", err)
		}
	} else if !strings.HasPrefix(newUrl, "http://") && !strings.HasPrefix(newUrl, "https://") {
		// Handle absolute URLs that do not start with "http://" or "https://"
		baseURL, err := url.Parse(fetchedUrl)
		if err != nil {
			return "", fmt.Errorf("error parsing URL: %w", err)
		}
		// Ensure there is a slash between the host and the new URL path
		if !strings.HasPrefix(newUrl, "/") {
			newUrl = "/" + newUrl
		}
		finalUrl, err = url.Parse(baseURL.Scheme + "://" + baseURL.Host + newUrl)
		if err != nil {
			return "", fmt.Errorf("error constructing final URL: %w", err)
		}
	} else {
		// Parse complete URL
		finalUrl, err = url.Parse(newUrl)
		if err != nil {
			return "", fmt.Errorf("error parsing URL: %w", err)
		}
	}

	// Remove specified unique parameters from the URL
	queryParams := finalUrl.Query()
	for _, param := range uniqueParameters {
		queryParams.Del(param)
	}
	finalUrl.RawQuery = queryParams.Encode()

	if removeFragment {
		finalUrl.Fragment = ""
	}

	// Return the modified URL
	return finalUrl.String(), nil
}
