package leveldb_test

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/ophum/humstack/pkg/store/leveldb"
)

func TestLevelDBStore(t *testing.T) {
	noti := make(chan string, 1000000)
	count := 0
	s, err := leveldb.NewLevelDBStore("./test", noti, false)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	wg := sync.WaitGroup{}
	for i := 0; i < 100000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			st := time.Now()
			s.Lock("test")
			log.Println(i, "locked! start")
			s.Put("test", struct {
				Test string
				Data string
			}{
				Test: "hoge",
				Data: fmt.Sprintf("fda%d", i),
			})
			count++
			en := time.Now()
			log.Println(i, "done unlock", st, en, en.Sub(st))
			s.Unlock("test")
		}(i)
	}
	wg.Wait()
	if count != 100000 {
		log.Fatalf("want: 100000\nexcepted: %d", count)
	}
}

func TestDeleteStore(t *testing.T) {
	noti := make(chan string, 1000000)
	s, err := leveldb.NewLevelDBStore("./test", noti, false)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	s.Put("test", struct {
		Data string
	}{
		Data: "hogehoge",
	})

	s.Delete("test")

	d := struct {
		Data string
	}{}
	err = s.Get("test", &d)
	if err != nil {
		if err.Error() != "Not Found" {
			t.Fatal("want: Not Found, expect: ", err.Error())
		}
	} else {
		t.Fatal("want: not exists, expect: exists")
	}
}
