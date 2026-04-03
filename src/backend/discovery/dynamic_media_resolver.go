package discovery

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
)

type MediaCandidate struct {
	URL          string
	LastModified time.Time
	Name         string
}

type DynamicOperatingSystemResolver struct {
	NetworkClient *http.Client
	UserAgent     string
}

func NewDynamicOperatingSystemResolver() *DynamicOperatingSystemResolver {
	return &DynamicOperatingSystemResolver{
		NetworkClient: &http.Client{
			Timeout: 30 * time.Second,
			CheckRedirect: func(
				req *http.Request,
				via []*http.Request,
			) error {
				return nil
			},
		},
		UserAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	}
}

func (resolver *DynamicOperatingSystemResolver) ResolveLatestArchitectureImage(
	sourceUrl string,
) (
	string,
	error,
) {
	lowerSourceUrl := strings.ToLower(
		sourceUrl,
	)

	if strings.HasSuffix(lowerSourceUrl, ".iso") ||
		strings.HasSuffix(lowerSourceUrl, ".img") ||
		strings.Contains(lowerSourceUrl, "sourceforge.net/projects") {
		return sourceUrl, nil
	}

	return resolver.scrapeDirectoryForLatestImage(
		sourceUrl,
	)
}

func (resolver *DynamicOperatingSystemResolver) scrapeDirectoryForLatestImage(
	directoryUrl string,
) (
	string,
	error,
) {
	request, err := http.NewRequest(
		"GET",
		directoryUrl,
		nil,
	)
	if err != nil {
		return "", err
	}
	
	request.Header.Set(
		"User-Agent",
		resolver.UserAgent,
	)

	response, err := resolver.NetworkClient.Do(
		request,
	)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf(
			"http_status_%d",
			response.StatusCode,
		)
	}

	bodyContentBytes, err := io.ReadAll(
		response.Body,
	)
	if err != nil {
		return "", err
	}
	
	bodyContent := string(
		bodyContentBytes,
	)

	// Regex to match hrefs and potential date information in common directory listings
	// Matches: href="filename.iso" and captures the full row for date parsing
	rowPattern := regexp.MustCompile(
		`(?i)<a[^>]+href=["']?([^"' >]+\.(iso|img))["' >]?[^>]*>(?:[^<]+)?</a>\s*([\d-]{4,10}\s+[\d:]{4,8}|[\d]{1,2}-\w{3}-\d{4}\s+[\d:]{4,8}|[\d]{1,2}\s+\w{3}\s+\d{4})?`,
	)
	
	allRows := rowPattern.FindAllStringSubmatch(
		bodyContent,
		-1,
	)

	if len(allRows) == 0 {
		return "", fmt.Errorf(
			"no_media_found_at_%s",
			directoryUrl,
		)
	}

	var candidates []MediaCandidate
	for _, row := range allRows {
		link := row[1]
		dateStr := row[3]
		lowerLink := strings.ToLower(
			link,
		)

		if strings.Contains(lowerLink, "netinst") ||
			strings.Contains(lowerLink, "minimal") ||
			strings.Contains(lowerLink, "mac") ||
			strings.Contains(lowerLink, "arm") ||
			strings.Contains(lowerLink, "zsync") ||
			strings.Contains(lowerLink, "sha256") ||
			strings.Contains(lowerLink, "md5") ||
			strings.Contains(lowerLink, "sig") {
			continue
		}

		parsedDate := time.Time{}
		if dateStr != "" {
			// Try common date formats
			formats := []string{
				"2006-01-02 15:04",
				"02-Jan-2006 15:04",
				"02 Jan 2006 15:04",
			}
			for _, fmtStr := range formats {
				if t, err := time.Parse(fmtStr, dateStr); err == nil {
					parsedDate = t
					break
				}
			}
		}

		candidates = append(
			candidates,
			MediaCandidate{
				URL:          link,
				LastModified: parsedDate,
				Name:         link,
			},
		)
	}

	if len(candidates) == 0 {
		return "", fmt.Errorf(
			"filtered_all_media_at_%s",
			directoryUrl,
		)
	}

	// Sort by date (descending), then by name (descending)
	sort.Slice(
		candidates,
		func(i, j int) bool {
			if !candidates[i].LastModified.Equal(candidates[j].LastModified) {
				return candidates[i].LastModified.After(candidates[j].LastModified)
			}
			return candidates[i].Name > candidates[j].Name
		},
	)

	selectedLink := candidates[0].URL

	parsedBase, err := url.Parse(
		response.Request.URL.String(),
	)
	if err != nil {
		return "", err
	}

	parsedLink, err := url.Parse(
		selectedLink,
	)
	if err != nil {
		return "", err
	}

	return parsedBase.ResolveReference(
		parsedLink,
	).String(), nil
}
