package leveldb

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"sync"

	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/syndtr/goleveldb/leveldb"
	leveldbErrors "github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type NoticeData struct {
	Key     string       `json:"key"`
	APIType meta.APIType `json:"apiType"`
	Before  string       `json:"before"`
	After   string       `json:"after"`
}

type LevelDBStore struct {
	db              *leveldb.DB
	lockTableLocker *sync.RWMutex
	lockTable       map[string]*sync.RWMutex
	notifier        chan string
	isDebug         bool
}

func NewLevelDBStore(dirPath string, notifier chan string, isDebug bool) (*LevelDBStore, error) {
	db, err := leveldb.OpenFile(filepath.Join(dirPath, "database.leveldb"), nil)
	if err != nil {
		return nil, err
	}

	return &LevelDBStore{
		db:              db,
		lockTableLocker: &sync.RWMutex{},
		lockTable:       map[string]*sync.RWMutex{},
		notifier:        notifier,
		isDebug:         isDebug,
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
	before, err := s.db.Get([]byte(key), nil)
	if err != nil {
		if err != leveldbErrors.ErrNotFound {
			return
		}
	}

	dataJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return
	}

	err = s.db.Put([]byte(key), dataJSON, nil)
	if err != nil {
		return
	}

	if s.isDebug {
		fmt.Println("=============== PUT  ==================")
		s.printDB()
	}

	obj := meta.Object{}
	if len(before) == 0 {
		if err := json.Unmarshal(dataJSON, &obj); err != nil {
			return
		}
	} else {
		if err := json.Unmarshal(before, &obj); err != nil {
			return
		}
	}
	noticeData := NoticeData{
		Key:     key,
		APIType: obj.Meta.APIType,
		Before:  string(before),
		After:   string(dataJSON),
	}

	noticeJSON, err := json.Marshal(noticeData)
	if err != nil {
		return
	}
	s.notifier <- string(noticeJSON)
}

func (s *LevelDBStore) Delete(key string) {
	before, err := s.db.Get([]byte(key), nil)
	if err != nil {
		if err != leveldbErrors.ErrNotFound {
			return
		}
	}
	obj := meta.Object{}
	if err := json.Unmarshal(before, &obj); err != nil {
		log.Println(err.Error())
		return
	}

	err = s.db.Delete([]byte(key), nil)
	if err != nil {
		log.Println(err.Error())
		return
	}

	if s.isDebug {
		fmt.Println("============== DELETE ================")
		s.printDB()
	}

	noticeJSON, err := json.Marshal(NoticeData{
		Key:     key,
		APIType: obj.Meta.APIType,
		Before:  string(before),
		After:   "",
	})
	if err != nil {
		log.Println(err.Error())
		return
	}

	s.notifier <- string(noticeJSON)
}

func (s *LevelDBStore) Lock(key string) {
	s.lockTableLocker.RLock()
	_, ok := s.lockTable[key]
	s.lockTableLocker.RUnlock()
	if !ok {
		s.lockTableLocker.Lock()
		s.lockTable[key] = &sync.RWMutex{}
		s.lockTableLocker.Unlock()
	}

	s.lockTableLocker.RLock()
	s.lockTable[key].Lock()
	s.lockTableLocker.RUnlock()
}

func (s *LevelDBStore) Unlock(key string) {
	s.lockTableLocker.RLock()
	defer s.lockTableLocker.RUnlock()
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
