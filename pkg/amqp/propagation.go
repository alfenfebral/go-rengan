package amqp

type AmqpHeadersCarrier map[string]interface{}

func (c AmqpHeadersCarrier) Get(key string) string {
	v, ok := c[key]
	if !ok {
		return ""
	}
	return v.(string)
}

func (c AmqpHeadersCarrier) Set(key string, value string) {
	c[key] = value
}

func (c AmqpHeadersCarrier) Keys() []string {
	i := 0
	r := make([]string, len(c))

	for k := range c {
		r[i] = k
		i++
	}

	return r
}
