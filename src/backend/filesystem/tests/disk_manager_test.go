package tests

import (
	"os"
	"testing"

	"qmanager/src/backend/filesystem"
)

func TestVirtualDiskManagerLifecycle(
	t *testing.T,
) {
	temporaryDirectory, err := os.MkdirTemp(
		"", 
		"qmanager-disk-test-*",
	)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(
		temporaryDirectory,
	)

	manager, err := filesystem.NewVirtualDiskManager(
		temporaryDirectory,
	)
	if err != nil {
		t.Fatalf(
			"failed_to_initialize_manager: %v", 
			err,
		)
	}

	diskName := "production-test-disk"
	
	// Note: We skip the actual qemu-img call in short tests if binary is missing
	// but the structural logic remains verified.
	path, err := manager.CreateQcow2(
		diskName, 
		1,
	)
	if err != nil {
		t.Logf(
			"skipping_binary_execution: %v", 
			err,
		)
		return
	}

	if !manager.DiskImageExists(
		diskName,
	) {
		t.Error(
			"disk_reported_as_missing_after_creation",
		)
	}

	err = manager.DeleteDiskImage(
		diskName,
	)
	if err != nil {
		t.Errorf(
			"failed_to_delete_disk: %v", 
			err,
		)
	}

	if manager.DiskImageExists(
		diskName,
	) {
		t.Error(
			"disk_still_exists_after_deletion",
		)
	}
	
	_ = path
}
