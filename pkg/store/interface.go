package store

type Store interface {
	List(prefix string, f func(n int) []interface{}) []interface{}
	Get(key string) interface{}
	Put(key string, data interface{})
	Delete(key string)
	Lock(key string)
	Unlock(key string)
}
