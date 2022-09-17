package querier

type Querier interface {
	Init(args ...interface{})
	Do(expression string) string
	Close()
}
