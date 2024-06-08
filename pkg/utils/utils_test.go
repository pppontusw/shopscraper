package utils

import (
	"testing"
	"time"
)

func TestEnsureFullUrl(t *testing.T) {
	// Test case 1: newUrl is empty
	newUrl := ""
	fetchedUrl := "https://example.com"
	expectedResult := ""
	result, err := EnsureFullUrl(newUrl, fetchedUrl, []string{}, false)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != expectedResult {
		t.Errorf("Expected %s, but got %s", expectedResult, result)
	}

	// Test case 2: newUrl is already a full URL
	newUrl = "https://example.com/product1"
	fetchedUrl = "https://example.com"
	expectedResult = "https://example.com/product1"
	result, err = EnsureFullUrl(newUrl, fetchedUrl, []string{}, false)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != expectedResult {
		t.Errorf("Expected %s, but got %s", expectedResult, result)
	}

	// Test case 3: newUrl is a relative path
	newUrl = "/product1"
	fetchedUrl = "https://example.com"
	expectedResult = "https://example.com/product1"
	result, err = EnsureFullUrl(newUrl, fetchedUrl, []string{}, false)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != expectedResult {
		t.Errorf("Expected %s, but got %s", expectedResult, result)
	}

	// Test case 4: newUrl is a relative path with query parameters
	newUrl = "/product1?param=value"
	fetchedUrl = "https://example.com"
	expectedResult = "https://example.com/product1?param=value"
	result, err = EnsureFullUrl(newUrl, fetchedUrl, []string{}, false)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != expectedResult {
		t.Errorf("Expected %s, but got %s", expectedResult, result)
	}

	// Test case 5: newUrl is a relative path with fragment identifier
	newUrl = "/product1#section"
	fetchedUrl = "https://example.com"
	expectedResult = "https://example.com/product1#section"
	result, err = EnsureFullUrl(newUrl, fetchedUrl, []string{}, false)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != expectedResult {
		t.Errorf("Expected %s, but got %s", expectedResult, result)
	}

	// Test case 6: newUrl contains ./
	newUrl = "./baz.php"
	fetchedUrl = "https://example.com/foo/bar.php"
	expectedResult = "https://example.com/foo/baz.php"
	result, err = EnsureFullUrl(newUrl, fetchedUrl, []string{}, false)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != expectedResult {
		t.Errorf("Expected %s, but got %s", expectedResult, result)
	}

	// Test case 7: newUrl contains &sid query parameter
	newUrl = "./baz.php?sid=12345"
	fetchedUrl = "https://example.com/foo/bar.php"
	expectedResult = "https://example.com/foo/baz.php"
	result, err = EnsureFullUrl(newUrl, fetchedUrl, []string{"sid"}, false)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != expectedResult {
		t.Errorf("Expected %s, but got %s", expectedResult, result)
	}

	// Test case 8: newUrl contains many query parameters, sid is stripped
	newUrl = "./baz.php?page=10&sid=12345&sort=asc"
	fetchedUrl = "https://example.com/foo/bar.php"
	expectedResult = "https://example.com/foo/baz.php?page=10&sort=asc"
	result, err = EnsureFullUrl(newUrl, fetchedUrl, []string{"sid"}, false)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != expectedResult {
		t.Errorf("Expected %s, but got %s", expectedResult, result)
	}

	// Test case 9: newUrl contains fragment that should be removed
	newUrl = "/bar#list=2345"
	fetchedUrl = "https://example.com/foo/bar"
	expectedResult = "https://example.com/bar"
	result, err = EnsureFullUrl(newUrl, fetchedUrl, []string{"sid"}, true)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != expectedResult {
		t.Errorf("Expected %s, but got %s", expectedResult, result)
	}

	// Test case 10: newUrl contains fragment that should not be removed
	newUrl = "/bar#list=2345"
	fetchedUrl = "https://example.com/foo/bar"
	expectedResult = "https://example.com/bar#list=2345"
	result, err = EnsureFullUrl(newUrl, fetchedUrl, []string{}, false)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != expectedResult {
		t.Errorf("Expected %s, but got %s", expectedResult, result)
	}
}

func TestGetPastTimeThreshold(t *testing.T) {
	tests := []struct {
		duration time.Duration
		expected time.Time
	}{
		{1 * time.Hour, time.Now().UTC().Add(-1 * time.Hour)},
		{5 * time.Hour, time.Now().UTC().Add(-5 * time.Hour)},
		{24 * time.Hour, time.Now().UTC().Add(-24 * time.Hour)},
		{-3 * time.Hour, time.Now().UTC().Add(-3 * time.Hour)},
	}

	for _, tt := range tests {
		t.Run(tt.expected.String(), func(t *testing.T) {
			got := GetPastTimeThreshold(tt.duration)
			// Allowing a margin for potential delays in computation
			if got.Before(tt.expected.Add(-5*time.Second)) || got.After(tt.expected.Add(5*time.Second)) {
				t.Errorf("GetPastTimeThreshold(%v) = %v, want %v", tt.duration, got, tt.expected)
			}
		})
	}
}
