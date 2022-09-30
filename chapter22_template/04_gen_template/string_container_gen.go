package gen

type StringContainer struct {
	s []string
}

func NewStringContainer() *StringContainer {
	return &StringContainer{s: []string{}}
}
func (c *StringContainer) Put(val string) {
	c.s = append(c.s, val)
}
func (c *StringContainer) Get() string {
	r := c.s[0]
	c.s = c.s[1:]
	return r
}
