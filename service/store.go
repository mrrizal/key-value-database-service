package service

import (
	"errors"
	"sync"
)

var ErrorNoSuchKey = errors.New("no such key")

type myMap struct {
	sync.RWMutex
	m map[string]string
}

type storeService struct {
	store myMap
}

func NewStoreService() *storeService {
	svc := storeService{
		myMap{m: make(map[string]string)},
	}
	return &svc
}

func (s *storeService) Put(key, value string) error {
	s.store.Lock()
	defer s.store.Unlock()

	s.store.m[key] = value
	return nil
}

func (s *storeService) Get(key string) (string, error) {
	s.store.RLock()
	defer s.store.RUnlock()

	if val, ok := s.store.m[key]; ok {
		return val, nil
	}
	return "", ErrorNoSuchKey
}

func (s *storeService) Delete(key string) error {
	delete(s.store.m, key)
	return nil
}
