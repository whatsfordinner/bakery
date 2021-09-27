package trace

type ContextCarrier struct {
	Fields map[string]string
}

func (c ContextCarrier) Get(key string) string {
	return c.Fields[key]
}

func (c ContextCarrier) Set(key string, value string) {
	c.Fields[key] = value
}

func (c ContextCarrier) Keys() []string {
	keys := make([]string, 0, len(c.Fields))
	for k := range c.Fields {
		keys = append(keys, k)
	}

	return keys
}
