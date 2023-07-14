package service

type MockStoreService struct {
	PutFunc    func(key, value string) error
	GetFunc    func(key string) (string, error)
	DeleteFunc func(key string) error
}

// Put mocks the Put method of the StoreService interface
func (m *MockStoreService) Put(key, value string) error {
	if m.PutFunc != nil {
		return m.PutFunc(key, value)
	}
	return nil
}

// Get mocks the Get method of the StoreService interface
func (m *MockStoreService) Get(key string) (string, error) {
	if m.GetFunc != nil {
		return m.GetFunc(key)
	}
	return "", nil
}

// Delete mocks the Delete method of the StoreService interface
func (m *MockStoreService) Delete(key string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(key)
	}
	return nil
}
