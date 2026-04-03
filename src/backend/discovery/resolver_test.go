package discovery

import (
	"testing"
)

func TestNewDynamicOperatingSystemResolver(
	t *testing.T,
) {
	resolver := NewDynamicOperatingSystemResolver()
	if resolver == nil {
		t.Fatal(
			"failed_to_initialize_resolver",
		)
	}

	if resolver.UserAgent == "" {
		t.Error(
			"resolver_missing_user_agent",
		)
	}
}
