package logger

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("FileTransactionLogger", func() {
	var (
		filename string
	)

	BeforeEach(func() {
		// Create a temporary file for testing
		file, err := ioutil.TempFile("", "transaction_log_*.txt")
		Expect(err).NotTo(HaveOccurred())
		defer file.Close()

		filename = file.Name()
	})

	AfterEach(func() {
		// Clean up the temporary file
		err := os.Remove(filename)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should create a new file transaction logger", func() {
		logger, err := NewFileTransactionLogger(filename)
		Expect(err).NotTo(HaveOccurred())
		defer logger.Close()

		Expect(logger).NotTo(BeNil())
	})
})

var _ = Describe("Put, Delete, Run", func() {
	var (
		filename          string
		transactionLogger *FileTransactionLogger
	)

	BeforeEach(func() {
		// Create a temporary file for testing
		file, err := os.OpenFile("./transaction_log_test.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
		Expect(err).NotTo(HaveOccurred())
		defer file.Close()

		filename = file.Name()

		// Create a new FileTransactionLogger for each test case
		transactionLogger, err = NewFileTransactionLogger(filename)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		// Clean up the temporary file
		err := os.Remove(filename)
		Expect(err).NotTo(HaveOccurred())

		// Close the transaction logger
		err = transactionLogger.Close()
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("WritePut", func() {
		It("should write a put event to the file", func() {
			key := "key1"
			value := "value1"

			transactionLogger.events = make(chan Event, 16)
			transactionLogger.WritePut(key, value)

			event := <-transactionLogger.events

			Expect(event.Type).To(Equal(EventPut))
			Expect(event.Key).To(Equal(key))
			Expect(event.Value).To(Equal(value))
		})
	})

	Describe("WriteDelete", func() {
		It("should write a delete event to the file", func() {
			key := "key2"

			transactionLogger.events = make(chan Event, 16)
			transactionLogger.WriteDelete(key)

			event := <-transactionLogger.events

			Expect(event.Type).To(Equal(EventDelete))
			Expect(event.Key).To(Equal(key))
			Expect(event.Value).To(BeEmpty())
		})

		It("should handle channel nil", func() {
			key := "key2"

			err := transactionLogger.WriteDelete(key)
			Expect(err).NotTo(BeNil())
		})
	})

	Describe("Run", func() {
		It("should write events to the file", func() {
			key := "key3"
			value := "value3"

			transactionLogger.Run()

			transactionLogger.WritePut(key, value)
			time.Sleep(50 * time.Millisecond)

			fileContent, err := ioutil.ReadFile(filename)
			Expect(err).NotTo(HaveOccurred())

			Expect(string(fileContent)).To(ContainSubstring(fmt.Sprintf("%d\t%d\t%s\t%s", 1, EventPut, key, value)))
		})
	})
})

var _ = Describe("Close", func() {
	var (
		filename          string
		transactionLogger *FileTransactionLogger
	)

	BeforeEach(func() {
		// Create a temporary file for testing
		file, err := os.OpenFile("./transaction_log_test.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
		Expect(err).NotTo(HaveOccurred())
		defer file.Close()

		filename = file.Name()

		// Create a new FileTransactionLogger for each test case
		transactionLogger, err = NewFileTransactionLogger(filename)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		// Clean up the temporary file
		err := os.Remove(filename)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Close", func() {
		It("should close the file", func() {
			err := transactionLogger.Close()
			Expect(err).NotTo(HaveOccurred())

			// Try to write to the closed file
			err = transactionLogger.WritePut("key4", "value4")
			Expect(err).To(HaveOccurred())
		})
	})
})

var _ = Describe("ReadEvents", func() {
	var (
		filename          string
		transactionLogger *FileTransactionLogger
	)

	BeforeEach(func() {
		// Create a temporary file for testing
		file, err := os.OpenFile("./transaction_log_test.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
		Expect(err).NotTo(HaveOccurred())
		defer file.Close()

		filename = file.Name()

		// Create a new FileTransactionLogger for each test case
		transactionLogger, err = NewFileTransactionLogger(filename)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		// Clean up the temporary file
		err := os.Remove(filename)
		Expect(err).NotTo(HaveOccurred())

		// Close the transaction logger
		err = transactionLogger.Close()
		Expect(err).NotTo(HaveOccurred())
	})

	It("should read events from the file", func() {
		eventsToWrite := []Event{
			{1, EventPut, "key5", "value5"},
			{2, EventDelete, "key6", ""},
		}

		// Write events to the file
		file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0644)
		Expect(err).NotTo(HaveOccurred())
		defer file.Close()

		for _, event := range eventsToWrite {
			_, err := fmt.Fprintf(file, "%d\t%d\t%s\t%s\n", event.Sequence, event.Type, event.Key, event.Value)
			Expect(err).NotTo(HaveOccurred())
		}

		eventChannel, errChannel := transactionLogger.ReadEvents()

		for _, expectedEvent := range eventsToWrite {
			Eventually(eventChannel).Should(Receive(Equal(expectedEvent)))
		}

		// Ensure there are no more events
		Consistently(eventChannel).ShouldNot(Receive())

		// Ensure there are no errors
		Eventually(errChannel).ShouldNot(Receive())
	})

	It("should handle parse errors", func() {
		// Write an invalid event to the file
		file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0644)
		Expect(err).NotTo(HaveOccurred())
		defer file.Close()

		_, err = fmt.Fprintln(file, "invalid_event")
		Expect(err).NotTo(HaveOccurred())

		_, errChannel := transactionLogger.ReadEvents()

		Eventually(errChannel).Should(Receive(MatchError(ContainSubstring("input parse error"))))
	})

	It("should handle out-of-sequence events", func() {
		// Write an out-of-sequence event to the file
		transactionLogger.lastSequence = 2
		file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0644)
		Expect(err).NotTo(HaveOccurred())
		defer file.Close()

		_, err = fmt.Fprintln(file, "2\t1\tkey7\tvalue7")
		Expect(err).NotTo(HaveOccurred())

		_, errChannel := transactionLogger.ReadEvents()

		Eventually(errChannel).Should(Receive(MatchError(ContainSubstring("transaction number out of sequence"))))
	})
})
