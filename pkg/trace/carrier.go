package trace

type ContextCarrier map[string]string

func (c ContextCarrier) Get(key string) string {
	return c[key]
}

func (c ContextCarrier) Set(key string, value string) {
	c[key] = value
}

func (c ContextCarrier) Keys() []string {
	keys := make([]string, 0, len(c))
	for k := range c {
		keys = append(keys, k)
	}

	return keys
}
