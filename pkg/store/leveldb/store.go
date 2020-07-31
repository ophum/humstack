package leveldb

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
	leveldbErrors "github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type LevelDBStore struct {
	db        *leveldb.DB
	lockTable map[string]*sync.RWMutex
}

func NewLevelDBStore(dirPath string) (*LevelDBStore, error) {
	db, err := leveldb.OpenFile(filepath.Join(dirPath, "database.leveldb"), nil)
	if err != nil {
		return nil, err
	}

	return &LevelDBStore{
		db:        db,
		lockTable: map[string]*sync.RWMutex{},
	}, nil
}

func (s *LevelDBStore) Close() error {
	return s.db.Close()
}

func (s *LevelDBStore) List(prefix string, f func(n int) []interface{}) error {
	iter := s.db.NewIterator(util.BytesPrefix([]byte(prefix)), nil)
	defer iter.Release()
	listJSON := [][]byte{}
	for iter.Next() {
		v := make([]byte, len(iter.Value()))
		copy(v, iter.Value())
		listJSON = append(listJSON, v)
	}

	err := iter.Error()
	if err != nil {
		return err
	}

	m := f(len(listJSON))
	for i, dataJSON := range listJSON {
		err := json.Unmarshal(dataJSON, m[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *LevelDBStore) Get(key string, v interface{}) error {
	dataJSON, err := s.db.Get([]byte(key), nil)
	if err != nil {
		if err == leveldbErrors.ErrNotFound {
			return errors.New("Not Found")
		}
		return err
	}

	return json.Unmarshal(dataJSON, v)
}

func (s *LevelDBStore) Put(key string, data interface{}) {
	dataJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return
	}

	err = s.db.Put([]byte(key), dataJSON, nil)
	if err != nil {
		return
	}

	fmt.Println("=============== PUT  ==================")
	s.printDB()
}

func (s *LevelDBStore) Delete(key string) {
	err := s.db.Delete([]byte(key), nil)
	if err != nil {
		return
	}

	fmt.Println("============== DELETE ================")
	s.printDB()
}

func (s *LevelDBStore) Lock(key string) {
	if _, ok := s.lockTable[key]; !ok {
		s.lockTable[key] = &sync.RWMutex{}
	}

	s.lockTable[key].Lock()
}

func (s *LevelDBStore) Unlock(key string) {
	if _, ok := s.lockTable[key]; ok {
		s.lockTable[key].Unlock()
	}
}

func (s *LevelDBStore) printDB() {
	iter := s.db.NewIterator(util.BytesPrefix([]byte("")), nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		fmt.Printf("%s ==>\n%s\n", string(key), string(value))
	}
	iter.Release()
}
