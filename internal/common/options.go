package common

import (
	"encoding/json"
	"ev_pub/internal/errors"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
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
	a.typeKey = gjson.GetBytes(data, typeKey).String()
	if a.typeKey == "" {
		return errors.New("missing type in options")
	}

	a.val = gjson.GetBytes(data, valKey).String()
	data, err := sjson.DeleteBytes(data, typeKey)
	if err != nil {
		return errors.Wrap(err, `error in parsing options`)
	}
	data, err = sjson.DeleteBytes(data, valKey)
	if err != nil {
		return errors.Wrap(err, `error in parsing options`)
	}
	var unmarshalling map[string]string
	err = json.Unmarshal(data, &unmarshalling)
	if err != nil {
		return errors.Wrap(err, `error in parsing options`)
	}
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
