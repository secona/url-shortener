package urlshortener

import (
	"fmt"
	"net/url"
	"regexp"
)

var slugRegex = regexp.MustCompile(`^[a-zA-Z0-9\-]+$`)

func parseURL(raw string) (*url.URL, error) {
	url, err := url.Parse(raw)

	if err != nil {
		return nil, err
	}

	// check if the scheme is either "http" or "https"
	if url.Scheme != "http" && url.Scheme != "https" {
		return nil, fmt.Errorf("Invalid URL: must be either \"http\" or \"https\".")
	}

	// check if host is present
	if url.Host == "" {
		return nil, fmt.Errorf("Invalid URL: empty host.")
	}

	return url, nil
}

func parseSlug(slug string) (string, error) {
	if !slugRegex.MatchString(slug) {
		return "", fmt.Errorf("Shortened link must only contain alphabets, numbers, and hyphens!")
	}

	return slug, nil
}
