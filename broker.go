package messagequeue

type Broker struct {
	Driver
}

func NewBroker() *Broker {
	return &Broker{}
}

type BrokerConfig struct {
	Driver string
	Config func(v interface{}) error `config:", lazyload"`
}

func NewBrokerConfig() *BrokerConfig {
	return &BrokerConfig{}
}
func (c *BrokerConfig) ApplyTo(b *Broker) error {
	d, err := NewDriver(c.Driver, c.Config)
	if err != nil {
		return err
	}
	b.Driver = d
	return nil
}
