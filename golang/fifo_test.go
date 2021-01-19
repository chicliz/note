package golang

import (
	"sort"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFIFO(t *testing.T) {
	f := NewFIFO()
	for i := 0; i < 10; i++ {
		f.Push(i)
	}
	f.Close()

	var errCount int32
	retArray := make([]int, 0)

	reader := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		reader.Add(1)
		go func() {
			defer reader.Done()
			ret, err := f.Pop()
			if err != nil {
				atomic.AddInt32(&errCount, 1)
			} else {
				retArray = append(retArray, ret.(int))
			}
		}()
	}
	reader.Wait()

	sort.Ints(retArray)
	for k, v := range retArray {
		assert.Equal(t, v, k)
	}
	assert.Equal(t, int32(0), errCount)
}

func TestFIFO_Closed(t *testing.T) {
	f := NewFIFO()

	retArray := make([]int, 0)
	var errCount int32

	reader := sync.WaitGroup{}
	for i := 0; i < 12; i++ {
		reader.Add(1)
		go func() {
			defer reader.Done()
			ret, err := f.Pop()
			if err != nil {
				atomic.AddInt32(&errCount, 1)
			} else {
				retArray = append(retArray, ret.(int))
			}
		}()
	}

	writer := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		writer.Add(1)
		go func(index int) {
			defer writer.Done()
			f.Push(index)
		}(i)
	}
	writer.Wait()
	f.Close()
	reader.Wait()

	sort.Ints(retArray)
	for k, v := range retArray {
		assert.Equal(t, v, k)
	}
	assert.Equal(t, int32(2), errCount)
}
