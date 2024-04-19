package utils

import (
	"net/url"
	"strings"
)

func EnsureFullUrl(newUrl, fetchedUrl string) string {
	var returnUrl string
	if strings.TrimSpace(newUrl) != "" {
		if !strings.HasPrefix(newUrl, "http") {
			baseURL, _ := url.Parse(fetchedUrl)
			returnUrl = baseURL.Scheme + "://" + baseURL.Host + newUrl
		} else {
			returnUrl = newUrl
		}
	}
	return returnUrl
}
