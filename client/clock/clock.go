package clock

type Clock struct {
	vector []int64
}

func NewClock(vector []int64) Clock {
	return Clock{vector: vector}
}

func (c *Clock) Increment(index int32) {
	c.vector[index] += 1
}

func (c *Clock) Update(vector []int64) {
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
