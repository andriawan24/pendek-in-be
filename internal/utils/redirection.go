package utils

import (
	"net/url"
	"strings"

	"github.com/medama-io/go-useragent"
	"github.com/mostafa-asg/ip2country"
)

func ParseDeviceType(ua useragent.UserAgent) string {
	if ua.IsMobile() {
		return "mobile"
	}

	if ua.IsTablet() {
		return "tablet"
	}

	return "desktop"
}

func ParseBrowser(ua useragent.UserAgent) string {
	if ua.Browser() == "" {
		return "unknown"
	}

	return ua.Browser().String()
}

func ParseTrafficSource(referrer string) string {
	if referrer == "" {
		return "direct"
	}

	parsedURL, err := url.Parse(referrer)
	if err != nil {
		return "unknown"
	}

	host := strings.ToLower(parsedURL.Host)
	host = strings.TrimPrefix(host, "www.")

	socialMedia := []string{"facebook.com", "twitter.com", "x.com", "instagram.com", "linkedin.com", "pinterest.com", "reddit.com", "tiktok.com", "youtube.com"}
	for _, domain := range socialMedia {
		if strings.Contains(host, domain) {
			return domain
		}
	}

	searchEngines := []string{"google.com", "bing.com", "yahoo.com", "duckduckgo.com"}
	for _, domain := range searchEngines {
		if strings.Contains(host, domain) {
			return domain
		}
	}

	return host
}

func ParseCountryFromIp(ipAddress string) string {
	country := ip2country.GetCountry(ipAddress)
	if country == "" {
		return "unknown"
	}
	return country
}
