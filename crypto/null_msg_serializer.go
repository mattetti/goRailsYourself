package crypto

import (
	"errors"
	"fmt"
	"reflect"
)

type NullMsgSerializer struct{}

func (s NullMsgSerializer) Serialize(vptr interface{}) (string, error) {
	return fmt.Sprint(vptr), nil
}

// Can only deserialize to a string.
func (s NullMsgSerializer) Unserialize(data string, vptr interface{}) error {
	typ := reflect.TypeOf(vptr)
	if typ.Kind() != reflect.Ptr {
		errors.New("You passed an interface which isn't a pointer")
	}
	v := reflect.ValueOf(vptr).Elem()
	v.SetString(data)
	return nil
}
