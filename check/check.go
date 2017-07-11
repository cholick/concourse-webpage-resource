package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/cholick/concourse-webpage-resource/models"
)

func main() {
	err := DoCheck(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func DoCheck(in io.Reader, out io.Writer) error {
	var request models.CheckRequest
	err := json.NewDecoder(in).Decode(&request)
	if err != nil {
		return errors.New(fmt.Sprintf("Error unmarshalling input: %v", err.Error()))
	}
	versions := []models.Version{}

	resp, err := http.Head(request.Source.Url)
	if err != nil {
		return errors.New(fmt.Sprintf("Error making HEAD request: %v", err.Error()))
	}
	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Non-OK response from server: %v", resp.StatusCode))
	}

	currentVersion := models.Version{}
	eTag := resp.Header[http.CanonicalHeaderKey("ETag")]
	if len(eTag) > 0 {
		currentVersion.ETag = eTag[0]
	}
	lastModified := resp.Header[http.CanonicalHeaderKey("Last-Modified")]
	if len(lastModified) > 0 {
		currentVersion.LastModified = lastModified[0]
	}
	if currentVersion.ETag == "" && currentVersion.LastModified == "" {
		//todo: use checksum of entire page?
		return errors.New("Resource requires etag or last-modified header")
	}

	if currentVersion != request.Version {
		versions = append(versions, currentVersion)
	}

	err = json.NewEncoder(out).Encode(versions)
	return err
}
