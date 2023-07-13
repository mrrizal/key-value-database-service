package service

import "errors"

var ErrorNoSuchKey = errors.New("no such key")

type storeService struct {
	store map[string]string
}

func NewStoreService() *storeService {
	svc := storeService{
		make(map[string]string),
	}
	return &svc
}

func (s *storeService) Put(key, value string) error {
	s.store[key] = value
	return nil
}

func (s *storeService) Get(key string) (string, error) {
	if val, ok := s.store[key]; ok {
		return val, nil
	}
	return "", ErrorNoSuchKey
}

func (s *storeService) Delete(key string) error {
	delete(s.store, key)
	return nil
}
