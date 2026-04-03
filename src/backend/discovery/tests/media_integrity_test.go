package tests

import (
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"qmanager/src/backend/discovery"
	"qmanager/src/core"
)

func TestOperatingSystemMediaIntegrity(
	t *testing.T,
) {
	if testing.Short() {
		t.Skip(
			"skipping_media_integrity_validation_in_short_mode",
		)
	}

	catalogConfigurationPath := "../../../../config/catalog"
	operatingSystemCatalog, err := core.LoadConfigurationCatalogFromDirectory(
		catalogConfigurationPath,
	)
	if err != nil {
		t.Fatalf(
			"failed_to_load_catalog_for_testing: %v",
			err,
		)
	}

	dynamicResolver := discovery.NewDynamicOperatingSystemResolver()
	httpClient := &http.Client{
		Timeout: 60 * time.Second,
	}

	for _, osMetadata := range operatingSystemCatalog.OperatingSystems {
		t.Run(
			osMetadata.ID,
			func(
				t *testing.T,
			) {
				resolvedMediaUrl, err := dynamicResolver.ResolveLatestArchitectureImage(
					osMetadata.ISOURL,
				)
				if err != nil {
					t.Fatalf(
						"failed_to_resolve_media_url_for_%s: %v",
						osMetadata.ID,
						err,
					)
				}

				headRequest, err := http.NewRequest(
					"HEAD",
					resolvedMediaUrl,
					nil,
				)
				if err != nil {
					t.Fatalf(
						"failed_to_create_head_request: %v",
						err,
					)
				}
				headRequest.Header.Set(
					"User-Agent",
					"QManager-CI-Validator/1.0",
				)

				headResponse, err := httpClient.Do(
					headRequest,
				)
				if err != nil {
					t.Fatalf(
						"head_transport_failure: %v",
						err,
					)
				}
				defer headResponse.Body.Close()

				if headResponse.StatusCode != http.StatusOK {
					t.Fatalf(
						"invalid_head_response_status_%d_for_%s",
						headResponse.StatusCode,
						resolvedMediaUrl,
					)
				}

				totalSizeBytes := headResponse.ContentLength
				if totalSizeBytes <= 0 {
					totalSizeBytes = 500 * 1024 * 1024 // Fallback to 500MB
				}

				// Download approximately 1% as requested
				verificationBytes := totalSizeBytes / 100
				if verificationBytes < 1024*1024 {
					verificationBytes = 1024 * 1024 // Min 1MB
				}
				if verificationBytes > 25*1024*1024 {
					verificationBytes = 25 * 1024 * 1024 // Max 25MB for CI stability
				}

				getRequest, err := http.NewRequest(
					"GET",
					resolvedMediaUrl,
					nil,
				)
				if err != nil {
					t.Fatalf(
						"failed_to_create_get_request: %v",
						err,
					)
				}
				getRequest.Header.Set(
					"User-Agent",
					"QManager-CI-Validator/1.0",
				)

				getResponse, err := httpClient.Do(
					getRequest,
				)
				if err != nil {
					t.Fatalf(
						"get_transport_failure: %v",
						err,
					)
				}
				defer getResponse.Body.Close()

				if getResponse.StatusCode != http.StatusOK {
					t.Fatalf(
						"invalid_get_response_status_%d_for_%s",
						getResponse.StatusCode,
						resolvedMediaUrl,
					)
				}

				limitReader := io.LimitReader(
					getResponse.Body,
					verificationBytes,
				)

				actualBytesRead, err := io.Copy(
					io.Discard,
					limitReader,
				)
				if err != nil {
					t.Fatalf(
						"failed_to_stream_media_content: %v",
						err,
					)
				}

				if actualBytesRead < 1024*1024 {
					t.Errorf(
						"resolved_media_is_too_small: read only %d bytes",
						actualBytesRead,
					)
				}

				fmt.Printf(
					"Validated %s -> %s (Verified %d bytes successfully)\n",
					osMetadata.ID,
					resolvedMediaUrl,
					actualBytesRead,
				)
			},
		)
	}
}
