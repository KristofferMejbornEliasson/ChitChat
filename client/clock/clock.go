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

func ensureArrayBigEnough(array *[]int64, index int32) (bigEnoughArray *[]int64) {
	if array == nil {
		arr := make([]int64, index+1)
		array = &arr
	}
	if len(*array) <= int(index) {
		arr := make([]int64, index+1)
		copy(arr, *array)
		array = &arr
	}
	return array
}

func (c *Clock) Increment(index int32) {
	c.lock.Lock()
	ensureArrayBigEnough(&c.vector, index)
	c.vector[index] += 1
	c.lock.Unlock()
}

func (c *Clock) IncrementAndCopy(index int32) (updatedCopy []int64) {
	c.lock.Lock()
	defer c.lock.Unlock()
	ensureArrayBigEnough(&c.vector, index)
	c.vector[index] += 1
	updatedCopy = make([]int64, len(c.vector))
	copy(updatedCopy, c.vector)
	return updatedCopy
}

func (c *Clock) Update(other []int64) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.vector == nil {
		c.vector = other
		return
	}
	if other == nil {
		return
	}

	if len(c.vector) < len(other) { // If other array is longer, we update *it*, then copy it into ourselves.
		for i := 0; i < len(other); i++ {
			if c.vector[i] < other[i] {
				other[i] = c.vector[i]
			}
		}
		c.vector = other
	} else if len(c.vector) == len(other) {
		for i := 0; i < len(c.vector); i++ {
			if c.vector[i] < other[i] {
				c.vector[i] = other[i]
			}
		}
	} else { // If other array is shorter, we update the elements we can
		for i := 0; i < len(other); i++ {
			if c.vector[i] < other[i] {
				c.vector[i] = other[i]
			}
		}
	}
}
