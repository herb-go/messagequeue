package messagequeue

import (
	"encoding/json"
)

//Option broker option interface.
type Option interface {
	ApplyTo(*Broker) error
}

// OptionConfigMap option config in map format.
type OptionConfigMap struct {
	Driver string
	Config ConfigMap
}

//ApplyTo apply option to message queue broker.
func (o *OptionConfigMap) ApplyTo(broker *Broker) error {
	driver, err := NewDriver(o.Driver, &o.Config, "")
	if err != nil {
		return err
	}
	broker.Driver = driver
	return nil
}

//NewOptionConfigMap create new option config.
func NewOptionConfigMap() *OptionConfigMap {
	return &OptionConfigMap{
		Config: map[string]interface{}{},
	}
}

// Config confit interface
type Config interface {
	//Get get value form given key.
	//Return any error if raised.
	Get(key string, v interface{}) error
}

//ConfigMap Map  format config
type ConfigMap map[string]interface{}

//Get get value form given key.
//Return any error if raised.
func (c *ConfigMap) Get(key string, v interface{}) error {
	i, ok := (*c)[key]
	if !ok {
		return nil
	}
	bs, err := json.Marshal(i)
	if err != nil {
		return err
	}
	return json.Unmarshal(bs, v)
}

//Set set value to given key.
//Return any error if raised.
func (c *ConfigMap) Set(key string, v interface{}) error {
	(*c)[key] = v
	return nil
}
