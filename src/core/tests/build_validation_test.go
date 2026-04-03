package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestFullApplicationCompilation(
	t *testing.T,
) {
	if testing.Short() {
		t.Skip(
			"skipping_full_build_validation_in_short_mode",
		)
	}

	tempBuildDirectory, err := os.MkdirTemp(
		"",
		"qmanager-build-check-*",
	)
	if err != nil {
		t.Fatalf(
			"failed_to_create_temporary_build_directory: %v",
			err,
		)
	}
	defer os.RemoveAll(
		tempBuildDirectory,
	)

	projectRoot, err := filepath.Abs(
		"../../..",
	)
	if err != nil {
		t.Fatalf(
			"failed_to_resolve_project_root: %v",
			err,
		)
	}

	t.Run(
		"ValidateCliBinaryCompilation",
		func(
			t *testing.T,
		) {
			outputBinaryPath := filepath.Join(
				tempBuildDirectory,
				"qmanager-cli-test",
			)
			
			compilationCommand := exec.Command(
				"go",
				"build",
				"-v",
				"-o",
				outputBinaryPath,
				filepath.Join(
					projectRoot,
					"cmd/qmanager-cli/main.go",
				),
			)
			compilationCommand.Dir = projectRoot

			outputLog, err := compilationCommand.CombinedOutput()
			if err != nil {
				t.Fatalf(
					"cli_compilation_failed: %v\nOutput:\n%s",
					err,
					string(outputLog),
				)
			}

			if _, err := os.Stat(outputBinaryPath); os.IsNotExist(err) {
				t.Error(
					"cli_binary_not_generated_after_successful_build_command",
				)
			}
		},
	)

	t.Run(
		"ValidateCorePackageCompilation",
		func(
			t *testing.T,
		) {
			compilationCheckCommand := exec.Command(
				"go",
				"build",
				"-v",
				"./src/backend/...",
				"./src/core/...",
			)
			compilationCheckCommand.Dir = projectRoot

			outputLog, err := compilationCheckCommand.CombinedOutput()
			if err != nil {
				t.Fatalf(
					"core_packages_compilation_failed: %v\nOutput:\n%s",
					err,
					string(outputLog),
				)
			}
		},
	)
}
