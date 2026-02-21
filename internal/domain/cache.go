package domain

type Cacher interface {
	Get(key string) ([]byte, error)
	Add(key string, val []byte) error
	Done()
}
