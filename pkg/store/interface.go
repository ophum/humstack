package store

type Store interface {
	List(prefix string, f func(n int) []interface{}) error
	Get(key string, v interface{}) error
	Put(key string, data interface{})
	Delete(key string)
	Lock(key string)
	Unlock(key string)
}
