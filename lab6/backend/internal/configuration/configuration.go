package configuration

import (
	"errors"
	"fmt"
	"reflect"
)

type IConfiguration interface {
	Get(vType interface{}) (interface{}, error)
	AddConfiguration(value interface{}) IConfiguration
}

type Configuration struct {
	configurations map[string]interface{}
}

func (c *Configuration) Get(vType interface{}) (interface{}, error) {
	t := reflect.TypeOf(vType)

	if t == nil {
		panic("can't accept interface type")
	}

	if t.Kind() != reflect.Pointer {
		panic("can't get non pointer type")
	}

	key := ""
	if t.Elem() == nil {
		key = t.Name()
	} else {
		key = t.Elem().Name()
	}

	value, exists := c.configurations[key]
	if !exists {
		return nil, errors.New(fmt.Sprintf("can't find configuration of type %s", t))
	}

	return value, nil
}

func (c *Configuration) AddConfiguration(value interface{}) IConfiguration {
	t := reflect.TypeOf(value)

	if t == nil {
		panic("can't accept interface type")
	}

	if t.Kind() != reflect.Pointer {
		panic("configuration is not of pointer type")
	}

	key := ""
	if t.Elem() == nil {
		key = t.Name()
	} else {
		key = t.Elem().Name()
	}

	if _, exists := c.configurations[key]; exists {
		panic(fmt.Sprintf("configuration with type %s already exists", key))
	}

	c.configurations[key] = value
	return c
}

func CreateConfiguration() IConfiguration {
	return &Configuration{
		configurations: make(map[string]interface{}),
	}
}
