package clock

import (
	"fmt"
	"sync"
)

type Clock struct {
	vector []int64
	lock   sync.Mutex
}

func NewClock(vector []int64) *Clock {
	return &Clock{vector: vector}
}

func (c *Clock) String() string {
	c.lock.Lock()
	defer c.lock.Unlock()
	return fmt.Sprint(c.vector)
}

func (c *Clock) CopyVector() (copiedVector []int64) {
	c.lock.Lock()
	defer c.lock.Unlock()
	copiedVector = make([]int64, len(c.vector))
	copy(copiedVector, c.vector)
	return copiedVector
}

func (c *Clock) Increment(index int32) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.vector[index] += 1
}

func (c *Clock) IncrementAndCopy(index int32) (updatedCopy []int64) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.vector[index] += 1
	updatedCopy = make([]int64, len(c.vector))
	copy(updatedCopy, c.vector)
	return updatedCopy
}

func (c *Clock) Update(vector []int64) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.vector == nil {
		c.vector = vector
		return
	}
	if len(c.vector) <= len(vector) {
		for i := 0; i < len(c.vector); i++ {
			if c.vector[i] < vector[i] {
				c.vector[i] = vector[i]
			}
		}
		for len(c.vector) < len(vector) {
			c.vector = append(c.vector, vector[len(c.vector)-1])
		}
	} else {
		for i := 0; i < len(vector); i++ {
			if c.vector[i] < vector[i] {
				c.vector[i] = vector[i]
			}
		}
	}
}
