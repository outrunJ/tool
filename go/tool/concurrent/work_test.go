package concurrent

import (
	"testing"
	"sync"
	"log"
	"time"
)

type TestTask struct {
	name string
}

func (t *TestTask) Task() {
	log.Println(t.name)
	time.Sleep(time.Second)
}

func TestWork(t *testing.T) {
	names := []string{
		"steve",
		"bob",
		"mary",
		"therese",
		"jason",
	}
	p := NewWork(2)

	var wg sync.WaitGroup
	wg.Add(5 * len(names))

	for i := 0; i < 5; i++ {
		for _, name := range names {
			task := TestTask{name: name}
			go func() {
				p.Run(&task)
				wg.Done()
			}()
		}
	}
	wg.Wait()
	p.Shutdown()
}
