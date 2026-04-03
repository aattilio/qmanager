package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"qmanager/src/backend/discovery"
)

func TestDynamicOperatingSystemResolverUnit(
	t *testing.T,
) {
	testServer := httptest.NewServer(
		http.HandlerFunc(
			func(
				writer http.ResponseWriter, 
				request *http.Request,
			) {
				html := `
					<html>
						<body>
							<a href="unrelated.txt">Not an ISO</a>
							<a href="ubuntu-24.04-netinst.iso">Netinst version</a>
							<a href="ubuntu-24.04-desktop-amd64.iso">The correct one</a>
							<a href="/root-relative.iso">Root relative</a>
							<a href="latest.iso.zsync">Metadata file</a>
						</body>
					</html>
				`
				fmt.Fprintln(
					writer, 
					html,
				)
			},
		),
	)
	defer testServer.Close()

	resolver := discovery.NewDynamicOperatingSystemResolver()
	
	t.Run(
		"ResolveFromDirectoryListing", 
		func(
			t *testing.T,
		) {
			resolved, err := resolver.ResolveLatestArchitectureImage(
				testServer.URL,
			)
			if err != nil {
				t.Fatalf(
					"failed_to_resolve: %v", 
					err,
				)
			}

			// Based on logic, it should pick ubuntu-24.04-desktop-amd64.iso 
			// because it's the last non-filtered candidate.
			expected := testServer.URL + "/ubuntu-24.04-desktop-amd64.iso"
			if resolved != expected {
				t.Errorf(
					"expected %s, got %s", 
					expected, 
					resolved,
				)
			}
		},
	)

	t.Run(
		"ResolveDirectLink", 
		func(
			t *testing.T,
		) {
			directLink := "https://example.com/system.iso"
			resolved, err := resolver.ResolveLatestArchitectureImage(
				directLink,
			)
			if err != nil {
				t.Fatal(err)
			}

			if resolved != directLink {
				t.Errorf(
					"expected %s, got %s", 
					directLink, 
					resolved,
				)
			}
		},
	)
}
