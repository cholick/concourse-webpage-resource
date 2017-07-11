package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/cholick/concourse-webpage-resource/models"
)

func main() {
	if len(os.Args) < 2 {
		println("usage: " + os.Args[0] + " <destination>")
		os.Exit(1)
	}
	destination := os.Args[1]

	err := os.MkdirAll(destination, 0755)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	err = DoIn(destination, os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func DoIn(destination string, in io.Reader, out io.Writer) error {
	var request models.CheckRequest
	err := json.NewDecoder(in).Decode(&request)
	if err != nil {
		return errors.New(fmt.Sprintf("Error unmarshalling input: %v", err.Error()))
	}

	filename := strings.TrimSpace(request.Source.Filename)
	if filename == "" {
		return errors.New("Source property filename is require")
	}

	resp, err := http.Get(request.Source.Url)
	if err != nil {
		return errors.New(fmt.Sprintf("Error making GET request to [%s]: %v", request.Source.Url, err.Error()))
	}
	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Non-OK response from server: %v", resp.StatusCode))
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New(fmt.Sprintf("Error making GET request to [%s]: %v", request.Source.Url, err.Error()))
	}

	err = ioutil.WriteFile(
		fmt.Sprintf("%v/%v", destination, filename), bodyBytes, 0755,
	)
	if err != nil {
		return err
	}

	err = json.NewEncoder(out).Encode(models.InResponse{
		Version: request.Version,
	})

	return err
}
