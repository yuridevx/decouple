package types

import (
	"fmt"
	"reflect"
)

type Function struct {
	Argument *Type
	Return   *Type

	ContextPresent bool
	ErrorPresent   bool
	Function       reflect.Value
}

func (f *Function) IsArgumentPresent() bool {
	return f.Argument != nil
}
func (f *Function) IsReturnPresent() bool {
	return f.Return != nil
}

func (f *Function) Call(ctx reflect.Value, req reflect.Value) (reflect.Value, error) {
	var callVals []reflect.Value

	switch {
	case f.ContextPresent && f.IsArgumentPresent():
		callVals = []reflect.Value{ctx, req}
	case f.ContextPresent && !f.IsArgumentPresent():
		callVals = []reflect.Value{ctx}
	case !f.ContextPresent && f.IsArgumentPresent():
		callVals = []reflect.Value{req}
	case !f.ContextPresent && !f.IsArgumentPresent():
		callVals = []reflect.Value{}
	}

	retVals := f.Function.Call(callVals)

	switch {
	case f.IsReturnPresent() && !f.ErrorPresent:
		return retVals[0], nil
	case f.IsReturnPresent() && f.ErrorPresent:
		return retVals[0], UnwrapError(retVals[1])
	case !f.IsReturnPresent() && f.ErrorPresent:
		return reflect.Value{}, UnwrapError(retVals[0])
	}

	return reflect.Value{}, nil
}

func (f *Function) ValidateRequest() error {
	if !f.IsArgumentPresent() {
		return fmt.Errorf("request function Struct argument have to be present")
	}
	if !f.IsReturnPresent() {
		return fmt.Errorf("request function Struct return have to be present")
	}
	return nil
}

func (f *Function) ValidateSubscribe() error {
	if !f.IsArgumentPresent() {
		return fmt.Errorf("subcribe function Struct argument have to be present")
	}
	return nil
}

func ParseFunction(fn interface{}) (*Function, error) {
	v := reflect.ValueOf(fn)
	t := v.Type()
	if t.Kind() != reflect.Func {
		return nil, fmt.Errorf("fn must be a function")
	}

	rf := &Function{
		Function: v,
	}

	numIn := t.NumIn()
	numOut := t.NumOut()

	switch {
	case numIn == 0:
	case numIn == 1:
		reqType, err := ParseType(t.In(0))
		if err != nil {
			return nil, err
		}
		rf.Argument = reqType
	case numIn == 2:
		if !IsContext(t.In(0)) {
			return nil, fmt.Errorf("first argument must be a context.Context, %s given", t.Name())
		}
		rf.ContextPresent = true

		reqType, err := ParseType(t.In(1))
		if err != nil {
			return nil, err
		}
		rf.Argument = reqType
	default:
		return nil, fmt.Errorf("only () or (req Type) or (ctx context.Context, req Type) parameters are supported")
	}

	switch numOut {
	case 0: //nothing
	case 1:
		if IsError(t.Out(0)) {
			rf.ErrorPresent = true
			break
		}
		pt, err := ParseType(t.Out(0))
		if err != nil {
			return nil, err
		}
		rf.Return = pt
	case 2:
		pt, err := ParseType(t.Out(0))
		if err != nil {
			return nil, err
		}
		rf.Return = pt

		if !IsError(t.Out(1)) {
			return nil, fmt.Errorf("only Void, Type or error or (Type, error) are supported as return types")
		}
		rf.ErrorPresent = true
	default:
		return nil, fmt.Errorf("only Void, Type or error or (Type, error) are supported as return types")
	}

	return rf, nil
}
