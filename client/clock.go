package main

type Clock struct {
	vector []int64
}

func (c *Clock) Update(vector []int64) {
	if c.vector == nil {
		c.vector = vector
		return
	}
	for i := 0; i < len(c.vector); i++ {
		if c.vector[i] < vector[i] {
			c.vector[i] = vector[i]
		}
	}
	for len(c.vector) < len(vector) {
		c.vector = append(c.vector, vector[len(c.vector)-1])
	}
}
