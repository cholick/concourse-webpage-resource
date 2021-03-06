package main_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cholick/concourse-webpage-resource/check"
)

var _ = Describe("Check Test", func() {
	Context("when executed", func() {
		var out *bytes.Buffer

		BeforeEach(func() {
			out = bytes.NewBuffer(nil)
		})

		It("should return initial value if none present", func() {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				header := w.Header()
				header["Etag"] = []string{"42"}
				header["Last-Modified"] = []string{"Fri, 05 May 2017 03:17:02 GMT"}
			})
			server := httptest.NewServer(handler)
			defer server.Close()

			request := fmt.Sprintf(`{"source": {"url": "%s"}}`, server.URL)
			in := strings.NewReader(request)

			var err error
			go func() {
				err = DoCheck(in, out)
			}()

			Eventually(func() string {
				return strings.TrimSpace(string(out.Bytes()))
			}).Should(Equal(`[{"eTag":"42","lastModified":"Fri, 05 May 2017 03:17:02 GMT"}]`))
			Expect(err).To(BeNil())
		})

		It("should return empty array when existing version", func() {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				header := w.Header()
				header["Etag"] = []string{"2cbf81b6c481c25454e72b114d0fee55"}
			})
			server := httptest.NewServer(handler)
			defer server.Close()

			request := fmt.Sprintf(`{
				"source": {"url": "%s"},
				"version": {
					"eTag": "2cbf81b6c481c25454e72b114d0fee55"
				}
			}`, server.URL)
			in := strings.NewReader(request)

			var err error
			go func() {
				err = DoCheck(in, out)
			}()

			Consistently(func() error { return err }).Should(BeNil())
			Eventually(func() string {
				return strings.TrimSpace(string(out.Bytes()))
			}).Should(Equal(`[]`))
			Expect(err).To(BeNil())
		})

		It("returns error on bad request", func() {
			request := fmt.Sprintf("totally not json")
			in := strings.NewReader(request)

			var err error
			go func() {
				err = DoCheck(in, out)
			}()

			Eventually(func() error {
				return err
			}).ShouldNot(BeNil())
		})

		It("returns error on problem making request", func() {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(503)
			})
			server := httptest.NewServer(handler)
			defer server.Close()

			request := fmt.Sprintf(`{"source": {"url": "%s"}}`, server.URL)
			in := strings.NewReader(request)

			var err error
			go func() {
				err = DoCheck(in, out)
			}()

			Eventually(func() error {
				return err
			}).ShouldNot(BeNil())
		})
	})
})
