package tool

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func RetryGet(channel *chan interface{}, duration *time.Duration, timeout *time.Duration, get func() interface{}) {
	if duration == nil {
		*duration = 2 * time.Second
	}

	go func() {
		stop := false
		if timeout != nil {
			Delay(func() {
				stop = true
			}, *timeout)
		}

		var res interface{}
		ticker := time.NewTicker(*duration)
		for range ticker.C {
			res = get()
			if stop || res != nil {
				*channel <- res
				ticker.Stop()
				break
			}
		}
	}()
}

func RetryDo(channel *chan interface{}, n int, duration *time.Duration, timeout *time.Duration, do func(interface{})) {
	if duration == nil {
		*duration = 2 * time.Second
	}

	go func() {
		doneN := 0
		running := true
		if timeout != nil {
			Delay(func() {
				running = false
			}, *timeout)
		}

		var res interface{}
		for running {
			select {
			case res = <-*channel:
				do(res)
				doneN++
				if n == doneN {
					running = false
				}
			default:
				time.Sleep(*duration)
			}
		}
	}()
}

type Retrier interface {
	Get(get func() interface{})
	Do(do func(interface{}))
}

type retrier struct {
	channel  *chan interface{}
	num      int
	duration *time.Duration
	timeout  *time.Duration
}

func GetRetrier(duration *time.Duration, timeout *time.Duration) Retrier {
	r := &retrier{
		duration: duration,
		timeout:  timeout,
	}
	r.Reset()
	return r
}

func (r *retrier) Get(get func() interface{}) {
	RetryGet(r.channel, r.duration, r.timeout, get)
	r.num++
}
func (r *retrier) Do(do func(interface{})) {
	RetryDo(r.channel, r.num, r.duration, r.timeout, do)
}
func (r *retrier) Reset() {
	c := make(chan interface{}, 1)
	r.channel = &c
	r.num = 0
}

type rwMutexMap struct {
	lockMap map[string]*sync.RWMutex
	mapLock *sync.Mutex
}

func RWMutexMap() *rwMutexMap {
	return &rwMutexMap{
		lockMap: map[string]*sync.RWMutex{},
		mapLock: &sync.Mutex{},
	}
}

func (m *rwMutexMap) ensureLock(key string) {
	if m.lockMap[key] == nil {
		m.mapLock.Lock()
		m.lockMap[key] = new(sync.RWMutex)
		m.mapLock.Unlock()
	}
}

func (m *rwMutexMap) RLock(key string) {
	m.ensureLock(key)
	l := m.lockMap[key]
	l.RLock()
}

func (m *rwMutexMap) RUnlock(key string) {
	m.ensureLock(key)
	l := m.lockMap[key]
	l.RUnlock()
}

func (m *rwMutexMap) Lock(key string) {
	m.ensureLock(key)
	l := m.lockMap[key]
	l.Lock()
}

func (m *rwMutexMap) Unlock(key string) {
	m.ensureLock(key)
	l := m.lockMap[key]
	l.Unlock()
}

func Cleanup(fn func()) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fn()
	}()
}

func Delay(fn func(), delay time.Duration) {
	go func() {
		time.Sleep(delay)
		fn()
	}()
}

// fn must suicide
func Timeout(fn func() error, timeout time.Duration) error {
	fnChannel := make(chan error, 1)
	timeoutChannel := make(chan error, 1)
	defer close(timeoutChannel)

	go func() {
		fnChannel <- fn()
	}()

	go func() {
		time.Sleep(timeout)
		timeoutChannel <- fmt.Errorf("timeout at %v", timeout)
	}()

	select {
	case err := <-fnChannel:
		close(fnChannel)
		return err
	case err := <-timeoutChannel:
		close(timeoutChannel)
		return err
	}
}

func Parallel(
	fn func(interface{}) error,
	tasks *[]interface{},
	limit int,
	timeout time.Duration,
) (error, *[]error) {
	// init
	var retErr error

	// task
	tasksLen := len(*tasks)
	errSlice := make([]error, tasksLen)
	if tasksLen == 0 {
		return nil, &errSlice
	}
	doTask := func(ind int) {
		task := (*tasks)[ind]
		err := fn(task)
		errSlice[ind] = err
		if retErr == nil && err != nil {
			retErr = fmt.Errorf("tasks encouner error")
		}
	}

	// lock
	lockChannel := make(chan time.Duration, limit)
	for i := 0; i < limit; i++ {
		lockChannel <- 0
	}
	taskQueue := NewQueue()
	for ind := range *tasks {
		taskQueue.Push(ind)
	}

	push := func() {
		lockChannel <- 0
	}

	pull := func(fn func()) {
		select {
		case <-lockChannel:
			go fn()
		}
	}

	check := func() bool {
		return !taskQueue.Empty()
	}

	// exec
	for check() {
		pull(func() {
			// avoid multi coroutines pop
			if check() {
				doTask(taskQueue.Pop().(int))
			}
			push()
		})
	}

	return retErr, &errSlice
}

func Interruptable() chan struct{} {
	done := make(chan struct{})
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		select {
		case <-c:
			signal.Stop(c)
			done <- struct{}{}
		}
	}()
	return done
}

func Concurrent(done chan struct{}, fn func(), num int, ccu int, qps int) {
	interval := time.Duration(1e9/qps) * time.Nanosecond
	don := make(chan struct{}, 2)
	go func() {
		<-done
		for i := 0; i < ccu; i++ {
			don <- struct{}{}
		}
	}()

	//
	tasks := make(chan struct{})
	go func() {
		var wg sync.WaitGroup
		wg.Add(num)
		for i := 0; i < num; i++ {
			tasks <- struct{}{}
			wg.Done()
			time.Sleep(interval)
		}
		wg.Wait()
		close(tasks)
	}()

	//
	var wg sync.WaitGroup
	wg.Add(ccu)
	for i := 0; i < ccu; i++ {
		go func() {
			defer wg.Done()
			for range tasks {
				select {
				case <-don:
					return
				default:
					fn()
				}
			}
		}()
	}
	wg.Wait()
}

// TODO: this is pseudo-code
type PoolConnector interface {
	Conn() (interface{}, error)
	Close(interface{}) error
	Ping(interface{}) error
}
type PoolConnection struct {
	Connection interface{}
	status     int // 0 idle, 1 using
}

type Pool struct {
	max            int
	min            int
	retryDuration  *time.Duration
	connector      PoolConnector
	connectTimeout *time.Duration
	status         int // 0 ok, 1 expanding
	pool           *ListSynchronized
	lockConnect    *sync.Mutex
}

func NewPool(conn PoolConnector, max int, min int, retryDuration *time.Duration, connectTimeout *time.Duration) (*Pool, error) {
	if min < 1 {
		return nil, fmt.Errorf("min should gte 1")
	}
	if max < min {
		return nil, fmt.Errorf("max should gte min")
	}
	if retryDuration == nil {
		*retryDuration = 2 * time.Second
	}
	if *retryDuration < time.Second {
		return nil, fmt.Errorf("retryDuration should gte 1 second")
	}
	if connectTimeout == nil {
		*connectTimeout = 2 * time.Second
	}
	if *connectTimeout < time.Second {
		return nil, fmt.Errorf("connectTimeout should gte 1 second")
	}
	if conn == nil {
		return nil, fmt.Errorf("conn should not be nil")
	}

	p := &Pool{
		max:           max,
		min:           min,
		retryDuration: retryDuration,
		connector:     conn,
		status:        0,
		pool:          NewListSynchronized(),
		lockConnect:   &sync.Mutex{},
	}
	p.expand(p.min)
	return p, nil
}

func (p *Pool) connect() interface{} {
	err := fmt.Errorf("")
	var c interface{}
	for err != nil {
		c, err = p.connector.Conn()
		time.Sleep(*p.retryDuration)
	}
	return c
}

func (p *Pool) expand(num int) {
	go func() {
		if p.status == 1 {
			return
		}
		p.status = 1
		currentNum := p.pool.Len()
		targetNum := currentNum + num
		if targetNum > p.max {
			targetNum = p.max
		}
		plus := targetNum - currentNum
		for i := 0; i < plus; i++ {
			connection := p.connect()
			p.pool.PushBack(&PoolConnection{Connection: connection})
		}
		p.status = 0
	}()
}

func (p *Pool) Connect() (*PoolConnection, error) {
	var poolConn *PoolConnection
	stop := false
	Timeout(func() error {
		stop = true
		return nil
	}, *p.connectTimeout)

	p.lockConnect.Lock()
	for poolConn == nil && !stop {
		for ele := p.pool.Front(); ele != nil; ele = ele.Next() {
			poolConn = ele.Value.(*PoolConnection)
			if poolConn.status == 0 {
				break
			}
		}
		p.expand(1)
		time.Sleep(2 * time.Second)
	}
	if poolConn != nil {
		poolConn.status = 1
	}
	p.lockConnect.Unlock()

	if poolConn == nil {
		return nil, fmt.Errorf("connect timeout")
	}
	return poolConn, nil
}

func (p *Pool) Close(conn *PoolConnection) {
	p.lockConnect.Lock()
	conn.status = 0
	p.lockConnect.Unlock()
}

