package discovery

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

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

	return resolver.scrapeDirectoryForImage(
		sourceUrl,
	)
}

func (resolver *DynamicOperatingSystemResolver) scrapeDirectoryForImage(
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

	isoLinkPattern := regexp.MustCompile(
		`href=["']?([^"' >]+\.(iso|img))(["' >]|$)`,
	)
	
	allMatches := isoLinkPattern.FindAllStringSubmatch(
		bodyContent,
		-1,
	)

	if len(allMatches) == 0 {
		return "", fmt.Errorf(
			"no_iso_found_at_%s",
			directoryUrl,
		)
	}

	var bestLink string
	maxScore := -1

	for _, match := range allMatches {
		link := match[1]
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

		score := 0
		if strings.Contains(lowerLink, "desktop") {
			score += 10
		}
		if strings.Contains(lowerLink, "live") {
			score += 5
		}
		if strings.Contains(lowerLink, "full") {
			score += 3
		}
		if strings.Contains(lowerLink, "x86_64") || strings.Contains(lowerLink, "amd64") {
			score += 2
		}

		if score > maxScore {
			maxScore = score
			bestLink = link
		}
	}

	if bestLink == "" {
		bestLink = allMatches[len(allMatches)-1][1]
	}

	parsedBase, err := url.Parse(
		response.Request.URL.String(),
	)
	if err != nil {
		return "", err
	}

	parsedLink, err := url.Parse(
		bestLink,
	)
	if err != nil {
		return "", err
	}

	return parsedBase.ResolveReference(
		parsedLink,
	).String(), nil
}
