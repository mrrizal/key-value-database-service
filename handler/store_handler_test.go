package handler

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gorilla/mux"
	"github.com/mrrizal/key-value-database/logger"
	"github.com/mrrizal/key-value-database/service"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DeleteHandler", func() {
	var (
		router       *mux.Router
		server       *httptest.Server
		storeHandler *StoreHandler
		mockLogger   *logger.MockTransactionLogger
	)

	BeforeEach(func() {
		router = mux.NewRouter()
		storeService := service.NewStoreService()
		mockLogger = &logger.MockTransactionLogger{}
		storeHandler = &StoreHandler{
			storeService,
			mockLogger,
		}

		router.HandleFunc("/store/{key}", storeHandler.Delete).Methods(http.MethodDelete)

		server = httptest.NewServer(router)
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("Delete", func() {
		It("should return status OK on successful deletion", func() {
			request, err := http.NewRequest(http.MethodDelete, server.URL+"/store/your-key", nil)
			Expect(err).NotTo(HaveOccurred())

			response, err := http.DefaultClient.Do(request)
			Expect(err).NotTo(HaveOccurred())
			defer response.Body.Close()

			Expect(response.StatusCode).To(Equal(http.StatusOK))
		})

		It("should return status Internal Server Error on deletion error", func() {
			mockService := &service.MockStoreService{
				DeleteFunc: func(key string) error {
					return errors.New("error")
				},
			}
			storeHandler.svc = mockService

			request, err := http.NewRequest(http.MethodDelete, server.URL+"/store/your-key", nil)
			Expect(err).NotTo(HaveOccurred())

			response, err := http.DefaultClient.Do(request)
			Expect(err).NotTo(HaveOccurred())
			defer response.Body.Close()

			Expect(response.StatusCode).To(Equal(http.StatusInternalServerError))
		})
	})
})

var _ = Describe("GetHandler", func() {
	var (
		router       *mux.Router
		server       *httptest.Server
		storeHandler *StoreHandler
		mockService  *service.MockStoreService
	)

	BeforeEach(func() {
		router = mux.NewRouter()
		mockService = &service.MockStoreService{}
		storeHandler = &StoreHandler{
			svc: mockService,
		}

		router.HandleFunc("/store/{key}", storeHandler.Get).Methods(http.MethodGet)

		server = httptest.NewServer(router)
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("Get", func() {
		Context("when key exists", func() {
			It("should return the value with status OK", func() {
				expectedValue := "some value"
				mockService.GetFunc = func(key string) (string, error) {
					return expectedValue, nil
				}

				request, err := http.NewRequest(http.MethodGet, server.URL+"/store/your-key", nil)
				Expect(err).NotTo(HaveOccurred())

				response, err := http.DefaultClient.Do(request)
				Expect(err).NotTo(HaveOccurred())
				defer response.Body.Close()

				resp, _ := ioutil.ReadAll(response.Body)
				Expect(response.StatusCode).To(Equal(http.StatusOK))
				Expect(string(resp)).To(ContainSubstring(expectedValue))
			})
		})

		Context("when key does not exist", func() {
			It("should return status Not Found", func() {
				mockService.GetFunc = func(key string) (string, error) {
					return "", service.ErrorNoSuchKey
				}

				request, err := http.NewRequest(http.MethodGet, server.URL+"/store/non-existent-key", nil)
				Expect(err).NotTo(HaveOccurred())

				response, err := http.DefaultClient.Do(request)
				Expect(err).NotTo(HaveOccurred())
				defer response.Body.Close()

				Expect(response.StatusCode).To(Equal(http.StatusNotFound))
			})
		})

		Context("when an error occurs during retrieval", func() {
			It("should return status Internal Server Error", func() {
				mockService.GetFunc = func(key string) (string, error) {
					return "", errors.New("some error")
				}

				request, err := http.NewRequest(http.MethodGet, server.URL+"/store/error-key", nil)
				Expect(err).NotTo(HaveOccurred())

				response, err := http.DefaultClient.Do(request)
				Expect(err).NotTo(HaveOccurred())
				defer response.Body.Close()

				Expect(response.StatusCode).To(Equal(http.StatusInternalServerError))
			})
		})
	})
})

var _ = Describe("PutHandler", func() {
	var (
		router       *mux.Router
		server       *httptest.Server
		storeHandler *StoreHandler
		mockService  *service.MockStoreService
		mockLogger   *logger.MockTransactionLogger
	)

	BeforeEach(func() {
		router = mux.NewRouter()
		mockService = &service.MockStoreService{}
		mockLogger = &logger.MockTransactionLogger{}
		storeHandler = &StoreHandler{
			svc:    mockService,
			logger: mockLogger,
		}

		router.HandleFunc("/store/{key}", storeHandler.Put).Methods(http.MethodPut)

		server = httptest.NewServer(router)
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("Put", func() {
		It("should return status Created on successful put", func() {
			mockService.PutFunc = func(key, value string) error {
				return nil
			}

			requestBody := "some value"
			request, err := http.NewRequest(http.MethodPut, server.URL+"/store/your-key",
				strings.NewReader(requestBody))
			Expect(err).NotTo(HaveOccurred())

			response, err := http.DefaultClient.Do(request)
			Expect(err).NotTo(HaveOccurred())
			defer response.Body.Close()

			Expect(response.StatusCode).To(Equal(http.StatusCreated))
		})

		It("should return status Internal Server Error on put error", func() {
			mockService.PutFunc = func(key, value string) error {
				return errors.New("some error")
			}

			requestBody := "some value"
			request, err := http.NewRequest(http.MethodPut, server.URL+"/store/your-key", strings.NewReader(requestBody))
			Expect(err).NotTo(HaveOccurred())

			response, err := http.DefaultClient.Do(request)
			Expect(err).NotTo(HaveOccurred())
			defer response.Body.Close()

			Expect(response.StatusCode).To(Equal(http.StatusInternalServerError))
		})

		It("should return status Internal Server Error on read request body error", func() {
			request, err := http.NewRequest(http.MethodPut, server.URL+"/store/your-key", &MockReader{})
			Expect(err).NotTo(HaveOccurred())

			response, err := http.DefaultClient.Do(request)
			Expect(err).To(HaveOccurred())
			Expect(response).To(BeNil())
		})
	})
})

// MockReader is a mock implementation of io.Reader that always returns an error
type MockReader struct{}

func (m *MockReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("some error")
}

func (m *MockReader) Close() error {
	return nil
}
