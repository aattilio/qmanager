package tests

import (
	"os"
	"path/filepath"
	"testing"

	"qmanager/src/core"
)

func TestCatalogLoader(
	t *testing.T,
) {
	temporaryDirectory, err := os.MkdirTemp(
		"", 
		"qmanager-catalog-test-*",
	)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(
		temporaryDirectory,
	)

	xmlContent := `
		<catalog>
			<os id="test-linux" name="Test Linux" version="1.0" family="linux" arch="x86_64">
				<iso_url>https://example.com/test.iso</iso_url>
				<min_ram_mb>2048</min_ram_mb>
				<min_vcpus>2</min_vcpus>
				<min_disk_gb>20</min_disk_gb>
			</os>
		</catalog>
	`
	
	err = os.WriteFile(
		filepath.Join(temporaryDirectory, "test.xml"), 
		[]byte(xmlContent), 
		0644,
	)
	if err != nil {
		t.Fatal(err)
	}

	catalog, err := core.LoadConfigurationCatalogFromDirectory(
		temporaryDirectory,
	)
	if err != nil {
		t.Fatalf(
			"failed_to_load_catalog: %v", 
			err,
		)
	}

	if len(catalog.OperatingSystems) != 1 {
		t.Fatalf(
			"expected 1 OS, got %d", 
			len(catalog.OperatingSystems),
		)
	}

	os := catalog.OperatingSystems[0]
	if os.ID != "test-linux" || os.Name != "Test Linux" {
		t.Errorf(
			"loaded_metadata_mismatch: %+v", 
			os,
		)
	}
}
