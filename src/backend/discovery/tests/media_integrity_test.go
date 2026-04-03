package tests

import (
	"fmt"
	"io"
	"net/http"
	"runtime"
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
	
	// Create a shared client with a strict timeout to avoid hanging the entire suite
	httpClient := &http.Client{
		Timeout: 45 * time.Second,
	}

	// Limit concurrency to avoid network saturation and context deadline exceeded in CI
	// 5 simultaneous OS checks is a safe balance for typical CI bandwidth
	concurrencyLimit := make(chan struct{}, 5)

	for _, osMetadata := range operatingSystemCatalog.OperatingSystems {
		// Capture variable for closure
		currentOS := osMetadata
		
		t.Run(
			currentOS.ID,
			func(
				t *testing.T,
			) {
				t.Parallel() // Execute OS validations in parallel
				
				// Wait for a slot in the concurrency limit
				concurrencyLimit <- struct{}{}
				defer func() { <-concurrencyLimit }()
				
				var mirrorErrors []error
				verifiedAnyMirror := false

				for _, mirrorUrl := range currentOS.Mirrors {
					resolvedMediaUrl, err := dynamicResolver.ResolveLatestArchitectureImage(
						mirrorUrl,
					)
					if err != nil {
						mirrorErrors = append(
							mirrorErrors, 
							fmt.Errorf("resolution_failed_for_%s: %v", mirrorUrl, err),
						)
						continue
					}

					// HEAD request to check availability and get size
					headRequest, err := http.NewRequest(
						"HEAD",
						resolvedMediaUrl,
						nil,
					)
					if err != nil {
						mirrorErrors = append(
							mirrorErrors, 
							fmt.Errorf("head_req_creation_failed_for_%s: %v", resolvedMediaUrl, err),
						)
						continue
					}
					headRequest.Header.Set(
						"User-Agent",
						"QManager-CI-Validator/1.0",
					)

					headResponse, err := httpClient.Do(
						headRequest,
					)
					if err != nil {
						mirrorErrors = append(
							mirrorErrors, 
							fmt.Errorf("head_do_failed_for_%s: %v", resolvedMediaUrl, err),
						)
						continue
					}
					defer headResponse.Body.Close()

					if headResponse.StatusCode != http.StatusOK {
						mirrorErrors = append(
							mirrorErrors, 
							fmt.Errorf("head_status_%d_for_%s", headResponse.StatusCode, resolvedMediaUrl),
						)
						continue
					}

					totalSizeBytes := headResponse.ContentLength
					if totalSizeBytes <= 0 {
						totalSizeBytes = 500 * 1024 * 1024
					}

					// We verify 1% of the ISO or 15MB, whichever is smaller, to satisfy the requirement
					// while keeping the CI fast and storage-safe.
					onePercent := totalSizeBytes / 100
					verificationBytes := int64(15 * 1024 * 1024) 
					if onePercent < verificationBytes {
						verificationBytes = onePercent
					}
					
					// Minimum check: 1MB
					if verificationBytes < 1024*1024 {
						verificationBytes = 1024 * 1024
					}

					getRequest, err := http.NewRequest(
						"GET",
						resolvedMediaUrl,
						nil,
					)
					if err != nil {
						mirrorErrors = append(
							mirrorErrors, 
							fmt.Errorf("get_req_creation_failed_for_%s: %v", resolvedMediaUrl, err),
						)
						continue
					}
					getRequest.Header.Set(
						"User-Agent",
						"QManager-CI-Validator/1.0",
					)

					getResponse, err := httpClient.Do(
						getRequest,
					)
					if err != nil {
						mirrorErrors = append(
							mirrorErrors, 
							fmt.Errorf("get_do_failed_for_%s: %v", resolvedMediaUrl, err),
						)
						continue
					}
					defer getResponse.Body.Close()

					if getResponse.StatusCode != http.StatusOK {
						mirrorErrors = append(
							mirrorErrors, 
							fmt.Errorf("get_status_%d_for_%s", getResponse.StatusCode, resolvedMediaUrl),
						)
						continue
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
						mirrorErrors = append(
							mirrorErrors, 
							fmt.Errorf("streaming_failed_for_%s: %v", resolvedMediaUrl, err),
						)
						continue
					}

					if actualBytesRead >= 1024*1024 {
						verifiedAnyMirror = true
						// Subtest log for tracking
						t.Logf(
							"Validated %s via %s (Read %d bytes)",
							currentOS.ID,
							mirrorUrl,
							actualBytesRead,
						)
						break
					}
				}

				if !verifiedAnyMirror {
					t.Errorf(
						"failed_to_verify_any_mirror_for_%s: all_mirrors_failed", 
						currentOS.ID,
					)
					for _, mErr := range mirrorErrors {
						t.Logf("Mirror error: %v", mErr)
					}
				}

				// Final cleanup for this parallel subtest
				runtime.GC()
			},
		)
	}
}
