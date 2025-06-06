package geturl

import (
	"log/slog"
	"net/http"
	"net/url"
	"strings"
)

func GetServiceURL(r *http.Request, routes map[string]string) (*url.URL, string, error) {
	path := r.URL.Path
	slog.Debug("Resolving service for path", slog.String("path", path))

	var matchedPrefix string

	for prefix := range routes {

		if strings.HasPrefix(path, prefix) {
			if len(prefix) > len(matchedPrefix) {
				matchedPrefix = prefix
				slog.Debug("Found matching prefix", slog.String("matchedPrefix", matchedPrefix))
			}
		}
	}

	if matchedPrefix == "" {
		slog.Error("Can't find matchedPrefix for path", slog.String("path", path))
		return nil, "", http.ErrAbortHandler
	}

	serviceBase := routes[matchedPrefix]

	parsedURL, err := url.Parse(serviceBase)
	if err != nil {
		slog.Error("Failed to parse service URL",
			slog.String("serviceBase", serviceBase),
			slog.Any("error", err))
	}
	return parsedURL, matchedPrefix, err
}
