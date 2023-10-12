package consumer

type SourceType int

const (
	SourceTypeRabbit SourceType = 1
	// SourceTypeHook   SourceType = 2
	// SourceTypeKafka  SourceType = 3
)

type Config[T any] struct {
	Source SourceType
}

func (c *Config[T]) CreateConsumer(conn interface{}) (*Consumer[T], error) {
	// if c.Source == SourceTypeRabbit {
	// 	rbConn
	// }
	return nil, nil
}

// func (c *Config[T]) CreateConsumerRetry(conn interface{},
//) (*ConsumerRetry[T], error) {
// 	return nil, nil
// }
