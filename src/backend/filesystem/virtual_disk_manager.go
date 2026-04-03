package filesystem

import (
	"fmt"
	"os"
	"os/exec"
)

type VirtualDiskManager struct {
	BaseStorageDirectory string
}

func NewVirtualDiskManager(storagePath string) (*VirtualDiskManager, error) {
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		return nil, fmt.Errorf("failed_to_create_storage_directory: %w", err)
	}
	return &VirtualDiskManager{BaseStorageDirectory: storagePath}, nil
}

func (manager *VirtualDiskManager) CreateQcow2(
	imageName string,
	sizeGigabytes int,
) (string, error) {
	diskImagePath := fmt.Sprintf(
		"%s/%s.qcow2",
		manager.BaseStorageDirectory,
		imageName,
	)
	
	command := exec.Command(
		"qemu-img",
		"create",
		"-f",
		"qcow2",
		diskImagePath,
		fmt.Sprintf("%dG", sizeGigabytes),
	)
	
	if err := command.Run(); err != nil {
		return "", fmt.Errorf("qemu_img_creation_failed: %w", err)
	}
	
	return diskImagePath, nil
}

func (manager *VirtualDiskManager) DeleteDiskImage(imageName string) error {
	diskImagePath := fmt.Sprintf(
		"%s/%s.qcow2",
		manager.BaseStorageDirectory,
		imageName,
	)
	return os.Remove(diskImagePath)
}

func (manager *VirtualDiskManager) DiskImageExists(imageName string) bool {
	diskImagePath := fmt.Sprintf(
		"%s/%s.qcow2",
		manager.BaseStorageDirectory,
		imageName,
	)
	_, err := os.Stat(diskImagePath)
	return !os.IsNotExist(err)
}
