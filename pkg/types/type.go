package types

import (
	"fmt"
	"reflect"
	"strings"
)

type Type struct {
	Name      string
	ValueType reflect.Type
}

func ParseTypeNamed(rt reflect.Type, name string) (*Type, error) {
	t := &Type{}
	switch rt.Kind() {
	case reflect.Struct:
		t.ValueType = rt
	default:
		return nil, fmt.Errorf("only Struct is supported, given %s", rt.String())
	}

	if name == "" {
		t.Name = t.ValueType.PkgPath() + "." + t.ValueType.Name()
	} else {
		t.Name = name
	}

	t.Name = strings.Replace(t.Name, "/", ".", -1)

	return t, nil
}

func ParseType(t reflect.Type) (*Type, error) {
	return ParseTypeNamed(t, "")
}
