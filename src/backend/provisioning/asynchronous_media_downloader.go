package provisioning

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type DownloadTask struct {
	URL      string
	Dest     string
	Progress func(current, total int64)
}

type byteCounter struct {
	total      int64
	downloaded int64
	onProgress func(current, total int64)
}

func (counter *byteCounter) Write(payload []byte) (int, error) {
	bytesRead := len(payload)
	counter.downloaded += int64(bytesRead)
	if counter.onProgress != nil {
		counter.onProgress(
			counter.downloaded,
			counter.total,
		)
	}
	return bytesRead, nil
}

func ExecuteDownload(task DownloadTask) error {
	outputFile, err := os.Create(task.Dest)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	response, err := http.Get(task.URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"download_failed_with_status_%d",
			response.StatusCode,
		)
	}

	progressCounter := &byteCounter{
		total:      response.ContentLength,
		onProgress: task.Progress,
	}
	
	_, err = io.Copy(
		outputFile,
		io.TeeReader(
			response.Body,
			progressCounter,
		),
	)
	
	return err
}
