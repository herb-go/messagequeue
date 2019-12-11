package messagequeue

//Option broker option interface.
type Option interface {
	ApplyTo(*Broker) error
}

// OptionConfig option config in map format.
type OptionConfig struct {
	Driver string
	Config func(interface{}) error `config:", lazyload"`
}

//ApplyTo apply option to message queue broker.
func (o *OptionConfig) ApplyTo(broker *Broker) error {
	driver, err := NewDriver(o.Driver, o.Config)
	if err != nil {
		return err
	}
	broker.Driver = driver
	return nil
}

//NewOptionConfig create new option config.
func NewOptionConfig() *OptionConfig {
	return &OptionConfig{}
}
