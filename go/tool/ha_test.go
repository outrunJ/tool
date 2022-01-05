package tool

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRWMutexMap(t *testing.T) {
	m := RWMutexMap()

	read := func(key string, i int) {
		fmt.Printf("%v , %v , read start\n", key, i)
		m.RLock(key)
		fmt.Printf("%v , %v , reading\n", key, i)
		time.Sleep(5 * time.Microsecond)
		m.RUnlock(key)
		fmt.Printf("%v , %v , read over\n", key, i)
	}

	write := func(key string, i int) {
		fmt.Printf("%v , %v , write start\n", key, i)
		m.Lock(key)
		fmt.Printf("%v , %v , writing\n", key, i)
		time.Sleep(5 * time.Microsecond)
		m.Unlock(key)
		fmt.Printf("%v , %v , write over\n", key, i)
	}

	go write("a", 1)
	go write("b", 1)
	go read("a", 2)
	go read("b", 2)
	go write("a", 3)
	go write("b", 3)

	time.Sleep(10 * time.Microsecond)
}

func TestParallel(t *testing.T) {
	assert := assert.New(t)

	err, errRst := Parallel(func(task interface{}) error {
		fmt.Println(task)
		if task.(int)%2 == 0 {
			return fmt.Errorf("a")
		} else {
			return nil
		}
	}, &[]interface{}{0, 1, 2, 3}, 2, time.Second)
	assert.Error(err)
	fmt.Println(errRst)
}
