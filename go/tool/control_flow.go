package tool

import (
	"container/list"
)

// pipeline
type PipelineWithDoneFnType func(interface{}) (bool, interface{}, error)
type PipelineWithDoneType struct {
	Done func() (interface{}, error)
	Do   func(PipelineWithDoneFnType) *PipelineWithDoneType
}

func PipelineReturnNotDone(b interface{}, c error) (bool, interface{}, error) {
	return false, b, c
}
func PipelineReturnDone(b interface{}, c error) (bool, interface{}, error) {
	return true, b, c
}
func PipelineWithDone(bag interface{}) *PipelineWithDoneType {
	fns := list.New()
	ptype := new(PipelineWithDoneType)

	ptype.Done = func() (interface{}, error) {
		var err error
		var done bool = false
		for ele := fns.Front(); !done && ele != nil; ele = ele.Next() {
			fn := ele.Value.(PipelineWithDoneFnType)
			if done, bag, err = fn(bag); err != nil {
				return bag, err
			}
		}

		return bag, nil
	}
	ptype.Do = func(f PipelineWithDoneFnType) *PipelineWithDoneType {
		fns.PushBack(f)
		return ptype
	}
	return ptype
}

type PipelineWithFnType func(interface{}) (interface{}, error)
type PipelineWithType struct {
	Done func() (interface{}, error)
	Do   func(PipelineWithFnType) *PipelineWithType
}

func PipelineWith(bag interface{}) *PipelineWithType {
	pipeline := PipelineWithDone(bag)
	ptype := new(PipelineWithType)
	ptype.Done = pipeline.Done

	ptype.Do = func(f PipelineWithFnType) *PipelineWithType {
		var err error
		_f := func(bag interface{}) (bool, interface{}, error) {
			bag, err = f(bag)
			return false, bag, err
		}
		pipeline.Do(_f)
		return ptype
	}
	return ptype
}

type PipelineFnType func() error
type PipelineType struct {
	Done func() error
	Do   func(PipelineFnType) *PipelineType
}

func Pipeline() *PipelineType {
	pipeline := PipelineWith(nil)

	ptype := new(PipelineType)
	ptype.Done = func() error {
		_, err := pipeline.Done()
		return err
	}

	ptype.Do = func(f PipelineFnType) *PipelineType {
		_f := func(_ interface{}) (interface{}, error) {
			return nil, f()
		}
		pipeline.Do(_f)
		return ptype
	}
	return ptype
}
