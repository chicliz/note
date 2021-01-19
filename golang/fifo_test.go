package golang

import (
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFIFO_Pop(t *testing.T) {
	f := NewFIFO()
	for i := 0; i < 10; i++ {
		f.Push(i)
	}
	f.Close()
	for i := 0; i < 10; i++ {
		ret, err := f.Pop()
		assert.Equal(t, i, ret)
		assert.Equal(t, nil, err)
	}
}

//并发pop
func TestFIFO(t *testing.T) {
	f := NewFIFO()
	for i := 0; i < 10; i++ {
		f.Push(i)
	}
	f.Close()
	checkPopItem(t, f, 10, 10)
}

//并发push and pop
func TestFIFO_Closed(t *testing.T) {
	f := NewFIFO()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		checkPopItem(t, f, 12, 10)
	}()

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
	wg.Wait()
}

func checkPopItem(t *testing.T, f *FIFO, popNum int, actualNum int) {
	var errCount int32
	retChan := make(chan int, 10)
	defer close(retChan)

	reader := sync.WaitGroup{}
	for i := 0; i < popNum; i++ {
		reader.Add(1)
		go func() {
			defer reader.Done()
			ret, err := f.Pop()
			fmt.Println(ret, err)
			if err != nil {
				atomic.AddInt32(&errCount, 1)
			} else {
				retChan <- ret.(int)
			}
		}()
	}
	reader.Wait()

	retArray := make([]int, 0)
	for i := 0; i < popNum-int(errCount); i++ {
		retArray = append(retArray, <-retChan)
	}
	sort.Ints(retArray)
	for k, v := range retArray {
		assert.Equal(t, v, k)
	}
	assert.Equal(t, int32(popNum-actualNum), errCount)
}
