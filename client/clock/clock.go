package clock

import (
	"fmt"
	"sync"
)

type Clock struct {
	vector []int64
	lock   sync.Mutex
}

// NewClock initialises a new vector clock structure with the given backing array.
func NewClock(vector []int64) *Clock {
	return &Clock{vector: vector}
}

func (c *Clock) String() string {
	c.lock.Lock()
	defer c.lock.Unlock()
	return fmt.Sprint(c.vector)
}

// Vector returns the backing array of this vector clock.
func (c *Clock) Vector() []int64 {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.vector == nil {
		c.vector = make([]int64, 0)
	}
	return c.vector
}

// grow ensures that the vector can contain at least 'length' number of elements.
func (c *Clock) grow(length int32) {
	if c.vector == nil {
		c.vector = make([]int64, length)
	} else if len(c.vector) < int(length) {
		arr := make([]int64, length)
		copy(arr, c.vector)
		c.vector = arr
	}
}

// Increment increments the logical clock at the given index.
func (c *Clock) Increment(index int32) {
	c.lock.Lock()
	c.grow(index + 1)
	c.vector[index]++
	c.lock.Unlock()
}

// IncrementAndCopy increments the logical clock at the given index, and returns the vector.
func (c *Clock) IncrementAndCopy(index int32) (updatedCopy []int64) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.grow(index + 1)
	c.vector[index]++
	return c.vector
}

// Update sets the receiver's array's elements to each be the maximum of the two.
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
		for i := 0; i < len(c.vector); i++ {
			if c.vector[i] > other[i] {
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
