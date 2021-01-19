package golang

import (
	"errors"
	"sync"
)

type FIFO struct {
	lock  sync.Mutex
	cond  sync.Cond
	queue []int

	exit  chan struct{}
}

//pop函数调用地方较多, 适用于多个消费者的模型
func (f *FIFO) Pop() (interface{}, error) {
	f.lock.Lock()
	defer f.lock.Unlock()
	for len(f.queue) == 0 {
		// When the queue is empty, invocation of Pop() is blocked until new item is enqueued.
		// When Close() is called, the f.closed is set and the condition is broadcasted.
		// Which causes this loop to continue and return from the Pop().
		if f.IsClosed() {
			return nil, errors.New("fifo closed")
		}

		f.cond.Wait()
	}
	id := f.queue[0]
	f.queue = f.queue[1:]
	return id, nil
}

func (f *FIFO) Push(s int) {
	f.lock.Lock()
	defer f.lock.Unlock()
	f.queue = append(f.queue, s)
	f.cond.Signal()
}

func (f *FIFO) IsClosed() bool {
	select {
	case <-f.exit:
		return true
	default:
		return false
	}
}

func (f *FIFO) Close() {
	close(f.exit)
	f.lock.Lock()
	defer f.lock.Unlock()
	f.cond.Broadcast()
}

func NewFIFO() *FIFO {
	f := &FIFO{
		queue: []int{},
	}
	f.cond.L = &f.lock
	f.exit = make(chan struct{})
	return f
}


