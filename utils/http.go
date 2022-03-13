package utils

import (
	"fmt"
	"io"
	"net/http"
)

func DownloadFile(url string) (io.ReadCloser, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Not able to download %s", url)
	}

	return response.Body, nil
}
