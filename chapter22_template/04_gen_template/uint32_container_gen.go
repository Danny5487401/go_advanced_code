package gen

type Uint32Container struct {
	s []uint32
}

func NewUint32Container() *Uint32Container {
	return &Uint32Container{s: []uint32{}}
}
func (c *Uint32Container) Put(val uint32) {
	c.s = append(c.s, val)
}
func (c *Uint32Container) Get() uint32 {
	r := c.s[0]
	c.s = c.s[1:]
	return r
}
