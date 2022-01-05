package tool

import (
	"container/list"
	"sync"
)

type Stack struct {
	list *list.List
}

func NewStack() *Stack {
	return &Stack{list: list.New()}
}

func (s *Stack) Push(value interface{}) {
	s.list.PushBack(value)
}

func (s *Stack) Pop() interface{} {
	ele := s.list.Back()
	if ele == nil {
		return nil
	}

	s.list.Remove(ele)
	return ele.Value
}

func (s *Stack) Peak() interface{} {
	ele := s.list.Back()
	if ele == nil {
		return nil
	}

	return ele.Value
}
func (s *Stack) Len() int {
	return s.list.Len()
}

func (s *Stack) Empty() bool {
	return s.list.Len() == 0
}

type Queue struct {
	list *list.List
}

func NewQueue() *Queue {
	return &Queue{list: list.New()}
}

func (s *Queue) Push(value interface{}) {
	s.list.PushBack(value)
}

func (s *Queue) Pop() interface{} {
	ele := s.list.Front()
	if ele == nil {
		return nil
	}

	s.list.Remove(ele)
	return ele.Value
}

func (s *Queue) Peak() interface{} {
	ele := s.list.Front()
	if ele == nil {
		return nil
	}

	return ele.Value
}
func (s *Queue) Len() int {
	return s.list.Len()
}

func (s *Queue) Empty() bool {
	return s.list.Len() == 0
}

type ListSynchronized struct {
	*list.List
	lock sync.Mutex
}

func NewListSynchronized() *ListSynchronized {
	return &ListSynchronized{
		List: list.New(),
		lock: sync.Mutex{},
	}
}

func (l *ListSynchronized) InsertAfter(v interface{}, mark *list.Element) *list.Element {
	l.lock.Lock()
	e := l.List.InsertAfter(v, mark)
	l.lock.Unlock()
	return e
}
func (l *ListSynchronized) InsertBefore(v interface{}, mark *list.Element) *list.Element {
	l.lock.Lock()
	e := l.List.InsertBefore(v, mark)
	l.lock.Unlock()
	return e
}
func (l *ListSynchronized) MoveAfter(e, mark *list.Element) {
	l.lock.Lock()
	l.List.MoveAfter(e, mark)
	l.lock.Unlock()
}
func (l *ListSynchronized) MoveBefore(e, mark *list.Element) {
	l.lock.Lock()
	l.List.MoveBefore(e, mark)
	l.lock.Unlock()
}
func (l *ListSynchronized) MoveToBack(e *list.Element) {
	l.lock.Lock()
	l.List.MoveToBack(e)
	l.lock.Unlock()
}
func (l *ListSynchronized) MoveToFront(e *list.Element) {
	l.lock.Lock()
	l.List.MoveToFront(e)
	l.lock.Unlock()
}
func (l *ListSynchronized) PushBack(v interface{}) {
	l.lock.Lock()
	l.List.PushBack(v)
	l.lock.Unlock()
}
func (l *ListSynchronized) PushBackList(other *list.List) {
	l.lock.Lock()
	l.List.PushBackList(other)
	l.lock.Unlock()
}
func (l *ListSynchronized) PushFront(v interface{}) {
	l.lock.Lock()
	l.List.PushFront(v)
	l.lock.Unlock()
}
func (l *ListSynchronized) PushFrontList(other *list.List) {
	l.lock.Lock()
	l.List.PushFrontList(other)
	l.lock.Unlock()
}
func (l *ListSynchronized) Remove(e *list.Element) {
	l.lock.Lock()
	l.List.Remove(e)
	l.lock.Unlock()
}
