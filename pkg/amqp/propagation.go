package pkg_amqp

type AmqpHeadersCarrier map[string]interface{}

func (carrier AmqpHeadersCarrier) Get(key string) string {
	v, ok := carrier[key]
	if !ok {
		return ""
	}
	return v.(string)
}

func (carrier AmqpHeadersCarrier) Set(key string, value string) {
	carrier[key] = value
}

func (carrier AmqpHeadersCarrier) Keys() []string {
	i := 0
	r := make([]string, len(carrier))

	for k := range carrier {
		r[i] = k
		i++
	}

	return r
}
