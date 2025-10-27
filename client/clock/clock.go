package clock

import "sync"

type Clock struct {
	vector []int64
	lock   sync.Mutex
}

func NewClock(vector []int64) *Clock {
	return &Clock{vector: vector}
}

func (c *Clock) Increment(index int32) {
	c.lock.Lock()
	c.vector[index] += 1
	c.lock.Unlock()
}

func (c *Clock) Update(vector []int64) {
	c.lock.Lock()
	if c.vector == nil {
		c.vector = vector
		c.lock.Unlock()
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
	c.lock.Unlock()
}
