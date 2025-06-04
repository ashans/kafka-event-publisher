package common

import (
	"encoding/json"
	"errors"
)

const (
	typeKey = `type`
	valKey  = `value`
)

type TypeValWithOptions struct {
	typeKey, val string
	options      map[string]string
}

func (a *TypeValWithOptions) UnmarshalJSON(data []byte) error {
	var unmarshalling map[string]string
	err := json.Unmarshal(data, &unmarshalling)
	if err != nil {
		return err
	}

	a.typeKey = unmarshalling[typeKey]
	if a.typeKey == "" {
		return errors.New("missing type")
	}
	a.val = unmarshalling[valKey]
	delete(unmarshalling, typeKey)
	delete(unmarshalling, valKey)
	a.options = unmarshalling
	return nil
}

func (a *TypeValWithOptions) Type() string {
	return a.typeKey
}
func (a *TypeValWithOptions) Val() string {
	return a.val
}

func (a *TypeValWithOptions) Options() map[string]string {
	return a.options
}
