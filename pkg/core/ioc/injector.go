package ioc

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"strconv"
	"sync"
)

type Injector struct {
	components sync.Map
}

func GetQualifiedKey(key string, qualifier string) string {
	if qualifier == "" {
		logrus.WithField("key", key).Warn("injector.GetQualifiedKey called with an empty qualifier")
		return key
	}
	return qualifier + "-" + key
}
func (c *Injector) InitModule(module string) {
	logrus.WithField("module", module).Info("injector.initmodule")
}

func (c *Injector) GetQualified(key string, qualifier string) interface{} {
	return c.Get(GetQualifiedKey(key, qualifier))
}
func (c *Injector) GetQualifiedInt(key string, qualifier string) int {
	return c.GetInt(GetQualifiedKey(key, qualifier))
}
func (c *Injector) GetQualifiedIntOrElse(key string, qualifier string, defaultValue int) int {
	return c.GetIntOrElse(GetQualifiedKey(key, qualifier), defaultValue)
}

func (c *Injector) GetQualifiedString(key string, qualifier string) string {
	return c.GetString(GetQualifiedKey(key, qualifier))
}
func (c *Injector) GetQualifiedStringOrElse(key string, qualifier string, defaultValue string) string {
	return c.GetStringOrElse(GetQualifiedKey(key, qualifier), defaultValue)
}

func (c *Injector) GetQualifiedBool(key string, qualifier string) bool {
	return c.GetBool(GetQualifiedKey(key, qualifier))
}
func (c *Injector) GetQualifiedBoolOrElse(key string, qualifier string, defaultValue bool) bool {
	return c.GetBoolOrElse(GetQualifiedKey(key, qualifier), defaultValue)
}

func (c *Injector) GetQualifiedOrElse(key string, qualifier string, defaultValue interface{}) interface{} {
	k := GetQualifiedKey(key, qualifier)
	res, ok := c.components.Load(k)
	if !ok {
		return defaultValue
	}
	return res
}

func (c *Injector) GetString(key string) string {
	logrus.WithField("key", key).Trace("injector.GetString")
	res, ok := c.components.Load(key)
	if !ok {
		panic(errors.New(fmt.Sprintf("Injector: missing component %s", key)))
	}

	val, err := c.AsString(res)
	if err != nil {
		panic(errors.New(fmt.Sprintf("Injector.GetString: invalid string value for %s ", key)))
	}
	return val
}

func (c *Injector) AsString(res interface{}) (string, error) {
	switch res.(type) {
	case string:
		return res.(string), nil
	}
	return "", errors.New(fmt.Sprintf("injector.AsString: invalid string type"))
}

func (c *Injector) GetInt(key string) int {
	logrus.WithField("key", key).Trace("injector.getInt")
	res, ok := c.components.Load(key)
	if !ok {
		panic(errors.New(fmt.Sprintf("Injector: missing component %s", key)))
	}

	val, err := c.AsInt(res)
	if err != nil {
		panic(errors.New(fmt.Sprintf("Injector.GetInt: invalid integer %s ", key)))
	}
	return val
}

func (c *Injector) AsInt(res interface{}) (int, error) {
	switch res.(type) {
	case int:
		return res.(int), nil
	case json.Number:
		int64, err := res.(json.Number).Int64()
		if err != nil {
			return 0, errors.New(fmt.Sprintf("injector.AsInt: fail to cast value %s from json.Number to int", res))
		}
		return int(int64), nil
	case string:
		val, err := strconv.Atoi(res.(string))
		if err != nil {
			return 0, errors.New(fmt.Sprintf("injector.AsInt: fail to cast value %s from string to integer ", res))
		}
		return val, nil
	}
	return 0, errors.New(fmt.Sprintf("injector.AsInt: invalid integer type"))
}

func (c *Injector) GetBool(key string) bool {
	logrus.WithField("key", key).Trace("injector.getBool")
	res, ok := c.components.Load(key)
	if !ok {
		panic(errors.New(fmt.Sprintf("Injector: missing component %s", key)))
	}

	val, err := c.AsBool(res)
	if err != nil {
		panic(errors.New(fmt.Sprintf("injector.GetBool: invalid bool %s ", key)))
	}
	return val
}

func (c *Injector) AsBool(res interface{}) (bool, error) {
	switch res.(type) {
	case bool:
		return res.(bool), nil
	case string:
		val, err := strconv.ParseBool(res.(string))
		if err != nil {
			return false, errors.New(fmt.Sprintf("injector.AsBool: fail to cast value %s from string to bool ", res))
		}
		return val, nil
	}
	return false, errors.New(fmt.Sprintf("injector.AsBool: invalid boolean type"))
}

func (c *Injector) GetStringSlice(key string) []string {
	logrus.WithField("key", key).Trace("injector.getStringSlice")
	res, ok := c.components.Load(key)
	if !ok {
		panic(errors.New(fmt.Sprintf("Injector: missing component %s", key)))
	}
	val, err := c.AsStringSlice(res)
	if err != nil {
		panic(errors.New(fmt.Sprintf("injector.GetStringSlice: invalid string %s ", key)))
	}
	return val
}

func (c *Injector) AsStringSlice(res interface{}) ([]string, error) {
	var result []string
	var errorList []error

	var handleValue = func(value interface{}) {
		strValue, err := c.AsString(value)
		if err != nil {
			errorList = append(errorList, err)
		}
		result = append(result, strValue)
	}

	switch res.(type) {
	case []string:
		for _, value := range res.([]string) {
			handleValue(value)
		}
	case []interface{}:
		for _, value := range res.([]interface{}) {
			handleValue(value)
		}
	default:
		return nil, errors.New(fmt.Sprintf("injector.AsStringSlice: invalid slice type"))
	}
	if len(errorList) != 0 {
		return nil, errors.New(fmt.Sprintf("injector.AsStringSlice contains value are not string: %s", errorList))
	}
	return result, nil
}

func (c *Injector) Get(key string) any {
	logrus.WithField("key", key).Trace("injector.get")
	res, ok := c.components.Load(key)
	if !ok {
		panic(errors.New(fmt.Sprintf("Injector: missing component %s", key)))
	}
	return res
}

func (c *Injector) GetIntOrElse(key string, defaultValue int) int {
	logrus.WithField("key", key).Trace("injector.getIntOrElse")
	res, ok := c.components.Load(key)
	if !ok {
		return defaultValue
	}
	val, err := c.AsInt(res)
	if err != nil {
		panic(errors.New(fmt.Sprintf("injector.GetIntOrElse: invalid int %s ", key)))
	}
	return val
}

func (c *Injector) GetBoolOrElse(key string, defaultValue bool) bool {
	logrus.WithField("key", key).Trace("injector.getBoolOrElse")
	res, ok := c.components.Load(key)
	if !ok {
		return defaultValue
	}
	val, err := c.AsBool(res)
	if err != nil {
		panic(errors.New(fmt.Sprintf("injector.GetBoolOrElse: invalid bool %s ", key)))
	}
	return val
}

func (c *Injector) GetStringOrElse(key string, defaultValue string) string {
	logrus.WithField("key", key).Trace("injector.getStringOrElse")
	res, ok := c.components.Load(key)
	if !ok {
		return defaultValue
	}
	val, err := c.AsString(res)
	if err != nil {
		panic(errors.New(fmt.Sprintf("injector.GetStringOrElse: invalid string %s ", key)))
	}
	return val
}

func (c *Injector) GetStringSliceOrElse(key string, defaultValue []string) []string {
	logrus.WithField("key", key).Trace("injector.getStringSliceOrElse")
	res, ok := c.components.Load(key)
	if !ok {
		return defaultValue
	}
	val, err := c.AsStringSlice(res)
	if err != nil {
		panic(errors.New(fmt.Sprintf("injector.GetStringSliceOrElse: invalid string %s ", key)))
	}
	return val
}

func (c *Injector) GetOrElse(key string, defaultValue interface{}) interface{} {
	logrus.WithField("key", key).Trace("injector.getorelse")
	res, ok := c.components.Load(key)
	if !ok {
		return defaultValue
	}
	return res
}

func (c *Injector) GetOrFunc(key string, defaultValueCreator func() interface{}) interface{} {
	logrus.WithField("key", key).Trace("injector.getorfunc")
	res, ok := c.components.Load(key)
	if !ok {
		return defaultValueCreator()
	}
	return res
}

func (c *Injector) Has(key string) bool {
	val, _ := c.components.Load(key)
	return val != nil
}

func (c *Injector) HasQualified(key string, qualifier string) bool {
	qualifiedKey := GetQualifiedKey(key, qualifier)
	val, _ := c.components.Load(qualifiedKey)
	return val != nil
}

func (c *Injector) RegisterQualified(key string, qualifier string, component interface{}) {
	logrus.WithField("key", key).WithField("qualifier", qualifier).Trace("injector.register")

	var qualifiedKey = GetQualifiedKey(key, qualifier)
	c.Register(qualifiedKey, component)
}

func (c *Injector) Register(key string, component interface{}) {
	logrus.WithField("key", key).Trace("injector.register")
	if c.Has(key) {
		panic(errors.New(fmt.Sprintf("Injector: component %s already registered, you must first unregister it explicitely", key)))
	}
	c.components.Store(key, component)
}

func (c *Injector) Replace(key string, component interface{}) {
	logrus.WithField("key", key).Trace("injector.replace")
	c.components.Store(key, component)
}

func (c *Injector) UnRegister(key string) {
	logrus.WithField("key", key).Trace("injector.unregister")
	if !c.Has(key) {
		panic(errors.New(fmt.Sprintf("Injector: component %s has not been registered", key)))
	}

	c.components.Delete(key)
}
