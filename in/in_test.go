package main_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cholick/concourse-webpage-resource/in"
)

var _ = Describe("In Test", func() {
	var tempDir string

	BeforeEach(func() {
		var err error
		tempDir, err = ioutil.TempDir("", "")
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		os.RemoveAll(tempDir)
	})

	It("writes out output", func() {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("some content"))
		})
		server := httptest.NewServer(handler)
		defer server.Close()

		request := fmt.Sprintf(`{"source": {"url": "%s","filename": "out.txt"}}`, server.URL)
		in := strings.NewReader(request)

		var err error
		go func() {
			err = DoIn(tempDir, in)
		}()

		Eventually(func() string {
			fileBytes, _ := ioutil.ReadFile(fmt.Sprintf("%v/out.txt", tempDir))
			return string(fileBytes)
		}).Should(Equal("some content"))
		Expect(err).To(BeNil())
	})

	It("returns error on non-200", func() {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(503)
		})
		server := httptest.NewServer(handler)
		defer server.Close()

		request := fmt.Sprintf(`{"source": {"url": "%s","filename": "out.txt"}}`, server.URL)
		in := strings.NewReader(request)

		var err error
		go func() {
			err = DoIn(tempDir, in)
		}()

		Eventually(func() error {
			return err
		}).ShouldNot(BeNil())
	})

	It("returns error if filename not specified", func() {
		in := strings.NewReader(`{"source": {"url": "https://www.example.com"}}`)

		var err error
		go func() {
			err = DoIn(tempDir, in)
		}()

		Eventually(func() error {
			return err
		}).ShouldNot(BeNil())
	})

	It("returns error on bad request", func() {
		request := fmt.Sprintf("totally not json")
		in := strings.NewReader(request)

		var err error
		go func() {
			err = DoIn(tempDir, in)
		}()

		Eventually(func() error {
			return err
		}).ShouldNot(BeNil())
	})
})
