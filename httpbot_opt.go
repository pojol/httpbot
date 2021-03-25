package httpbot

// Parm httpbot parm
type Parm struct {
	Name string

	// PrintReprot 是否输出报告
	PrintReprot bool
}

// Option consul discover config wrapper
type Option func(*Parm)

// WithPrintReprot 是否输出过程报告
func WithPrintReprot(report bool) Option {
	return func(c *Parm) {
		c.PrintReprot = report
	}
}

// WithName 传入bot的名称
func WithName(name string) Option {
	return func(c *Parm) {
		c.Name = name
	}
}
