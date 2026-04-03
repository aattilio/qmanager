package tests

import (
	"os/exec"
	"strings"
	"testing"
)

func TestQemuSystemAvailability(
	t *testing.T,
) {
	qemuBinary := "qemu-system-x86_64"
	
	command := exec.Command(
		qemuBinary,
		"--version",
	)
	
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf(
			"qemu_not_found_or_unresponsive: %v",
			err,
		)
	}

	versionString := string(output)
	if !strings.Contains(versionString, "QEMU emulator version") {
		t.Errorf(
			"unexpected_qemu_version_output: %s",
			versionString,
		)
	}
}

func TestQemuMinimalInvocation(
	t *testing.T,
) {
	command := exec.Command(
		"qemu-system-x86_64",
		"-help",
	)
	
	err := command.Run()
	if err != nil {
		t.Fatalf(
			"failed_to_invoke_qemu_help: %v",
			err,
		)
	}
}
