package main_test

import (
	"a21hc3NpZ25tZW50/service"
	"bytes"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("AIService", func() {
	var (
		mockClient *MockClient
		aiService  *service.AIService
	)

	BeforeEach(func() {
		mockClient = &MockClient{}
		aiService = &service.AIService{Client: mockClient}
	})

	Describe("ChatWithAI", func() {
		It("should return the correct response for a valid request", func() {
			mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
				// Mock respons API yang valid
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`[{"generated_text":"response"}]`)),
				}, nil
			}

			context := "context"
			query := "query"
			token := "token"

			// Memanggil fungsi ChatWithAI
			result, err := aiService.ChatWithAI(context, query, token)

			// Verifikasi hasilnya
			Expect(err).ToNot(HaveOccurred())
			Expect(result.GeneratedText).To(Equal("response"))
		})

		It("should return an error for an error response", func() {
			mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
				// Mock respons API yang error
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"error":"internal error"}`)),
				}, nil
			}

			context := "context"
			query := "query"
			token := "token"

			// Memanggil fungsi ChatWithAI
			result, err := aiService.ChatWithAI(context, query, token)

			// Verifikasi error yang dihasilkan
			Expect(err).To(HaveOccurred())
			Expect(result.GeneratedText).To(BeEmpty())
		})
	})
})

// MockClient untuk mock HTTP request pada pengujian
type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}
