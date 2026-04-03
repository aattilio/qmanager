package tests

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"qmanager/src/backend/provisioning"
)

func TestAsynchronousMediaDownloader(
	t *testing.T,
) {
	mockData := []byte(
		"fake-iso-binary-content-12345",
	)
	
	testServer := httptest.NewServer(
		http.HandlerFunc(
			func(
				writer http.ResponseWriter, 
				request *http.Request,
			) {
				writer.Write(
					mockData,
				)
			},
		),
	)
	defer testServer.Close()

	temporaryDirectory, err := os.MkdirTemp(
		"", 
		"qmanager-download-test-*",
	)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(
		temporaryDirectory,
	)

	destinationPath := filepath.Join(
		temporaryDirectory, 
		"test.iso",
	)
	
	progressCalled := false
	task := provisioning.DownloadTask{
		URL:  testServer.URL,
		Dest: destinationPath,
		Progress: func(
			current, 
			total int64,
		) {
			progressCalled = true
		},
	}

	err = provisioning.ExecuteDownload(
		task,
	)
	if err != nil {
		t.Fatalf(
			"download_failed: %v", 
			err,
		)
	}

	downloadedContent, err := os.ReadFile(
		destinationPath,
	)
	if err != nil {
		t.Fatal(err)
	}

	if string(downloadedContent) != string(mockData) {
		t.Error(
			"downloaded_content_mismatch",
		)
	}

	if !progressCalled {
		t.Error(
			"progress_callback_never_invoked",
		)
	}
}
