package cursors

import (
	"net"
	"net/url"
	"strings"

	"encore.dev/rlog"
)

// validateURL validates the URL to ensure it can be used for fetching and subscribing to cursor updates.
//
// Returns true if the url is valid, false otherwise.
func validateURL(targetURL string) bool {
	url, err := url.Parse(targetURL)
	if err != nil {
		rlog.Debug("failed to parse url", "err", err)
		return false
	}

	if cfg.AllowLocalhost() && (url.Host == "localhost" || url.Host == "127.0.0.1") {
		return true
	}

	if url.Host == "" {
		return false
	}
	if url.Fragment != "" {
		return false
	}
	if strings.HasSuffix(url.Host, "localhost") {
		return false
	}
	if ip := net.ParseIP(url.Host); ip != nil && (ip.IsLoopback() || ip.IsPrivate()) {
		return false
	}

	return true
}
