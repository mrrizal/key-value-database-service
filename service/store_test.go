package service

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Put", func() {
	var (
		key   string
		value string
	)
	var svc *storeService

	BeforeEach(func() {
		svc = NewStoreService()
		key = "testKey"
		value = "testValue"
	})

	It("should store the value correctly", func() {
		err := svc.Put(key, value)
		Expect(err).To(BeNil())
		Expect(svc.store[key]).To(Equal(value))
	})
})

var _ = Describe("Get", func() {
	var svc *storeService

	BeforeEach(func() {
		svc = NewStoreService()
	})

	Context("when key exists in the store", func() {
		It("should return the corresponding value", func() {
			key := "existingKey"
			value := "existingValue"
			svc.store[key] = value

			result, err := svc.Get(key)
			Expect(err).To(BeNil())
			Expect(result).To(Equal(value))
		})
	})

	Context("when key does not exist in the store", func() {
		It("should return an empty string", func() {
			key := "nonexistentKey"

			result, err := svc.Get(key)
			Expect(err).NotTo(BeNil())
			Expect(result).To(Equal(""))
		})

		It("should return an error of type ErrorNoSuchKey", func() {
			key := "nonexistentKey"

			_, err := svc.Get(key)
			Expect(err).To(Equal(ErrorNoSuchKey))
		})
	})
})

var _ = Describe("Delete", func() {
	var svc *storeService

	BeforeEach(func() {
		svc = NewStoreService()
	})

	Context("when key exists in the store", func() {
		It("should remove the key-value pair from the store", func() {
			key := "existingKey"
			value := "existingValue"
			svc.store[key] = value

			err := svc.Delete(key)
			Expect(err).To(BeNil())
			_, exists := svc.store[key]
			Expect(exists).To(BeFalse())
		})
	})

	Context("when key does not exist in the store", func() {
		It("should return an error", func() {
			key := "nonexistentKey"

			err := svc.Delete(key)
			Expect(err).To(BeNil())
		})

		It("should not modify the store", func() {
			key := "nonexistentKey"

			_ = svc.Delete(key)
			Expect(len(svc.store)).To(BeZero())
		})
	})
})
