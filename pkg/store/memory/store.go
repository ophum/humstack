package memory

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type MemoryStore struct {
	data      map[string]interface{}
	lockTable map[string]*sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data:      map[string]interface{}{},
		lockTable: map[string]*sync.RWMutex{},
	}
}

func (s *MemoryStore) List(prefix string, f func(n int) []interface{}) []interface{} {
	list := []interface{}{}
	for k, obj := range s.data {
		if strings.HasPrefix(k, prefix) {
			list = append(list, obj)
		}
	}

	m := f(len(list))
	for i, o := range list {
		a := reflect.Indirect(reflect.ValueOf(m[i]))
		a.Set(reflect.ValueOf(o))
	}

	return nil
}

func (s *MemoryStore) Get(key string) interface{} {
	if d, ok := s.data[key]; ok {
		return d
	}
	return nil
}

func (s *MemoryStore) Put(key string, data interface{}) {
	s.data[key] = data
	buf, _ := json.MarshalIndent(s.data, "", "  ")
	fmt.Println("============= PUT DATA ==============")
	fmt.Println(string(buf))
}

func (s *MemoryStore) Delete(key string) {
	delete(s.data, key)
	buf, _ := json.MarshalIndent(s.data, "", "  ")
	fmt.Println("============= DEL DATA ==============")
	fmt.Println(string(buf))
}

func (s *MemoryStore) Lock(key string) {
	if _, ok := s.lockTable[key]; !ok {
		s.lockTable[key] = &sync.RWMutex{}
	}

	s.lockTable[key].Lock()
}

func (s *MemoryStore) Unlock(key string) {
	if _, ok := s.lockTable[key]; ok {
		s.lockTable[key].Unlock()
	}
}
