package service

type StoreService interface {
	Put(key, value string) error
	Get(key string) (string, error)
	Delete(key string) error
}
