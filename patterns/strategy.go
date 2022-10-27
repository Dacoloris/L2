package pattern

type Context struct {
	strategy func()
}

func (c *Context) Execute() {
	c.strategy()
}

func (c *Context) SetStrategy(strategy func()) {
	c.strategy = strategy
}
