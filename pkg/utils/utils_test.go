package utils

import (
	"testing"
)

func TestEnsureFullUrl(t *testing.T) {
	// Test case 1: newUrl is empty
	newUrl := ""
	fetchedUrl := "https://example.com"
	expectedResult := ""
	result, err := EnsureFullUrl(newUrl, fetchedUrl)
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
	result, err = EnsureFullUrl(newUrl, fetchedUrl)
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
	result, err = EnsureFullUrl(newUrl, fetchedUrl)
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
	result, err = EnsureFullUrl(newUrl, fetchedUrl)
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
	result, err = EnsureFullUrl(newUrl, fetchedUrl)
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
	result, err = EnsureFullUrl(newUrl, fetchedUrl)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != expectedResult {
		t.Errorf("Expected %s, but got %s", expectedResult, result)
	}
}
