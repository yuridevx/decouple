package local

import (
	"decouple"
	"decouple/types"
	"fmt"
	"github.com/jinzhu/copier"
	"go.uber.org/multierr"
	"log"
	"reflect"
)

type Engine struct {
	Container *decouple.Container
}

func (e *Engine) Request(request interface{}, options ...decouple.CallOption) (interface{}, error) {
	reqV := reflect.ValueOf(request)
	reqT, err := types.ParseType(reqV.Type())
	if err != nil {
		panic(err)
	}
	reqFn := e.Container.GetRequestFunction(reqT.ValueType)
	if reqFn == nil {
		panic(fmt.Errorf("request handler not found for %s", reqT.Name))
	}

	opts := decouple.NewCallOptions(options)
	ctxV := reflect.ValueOf(opts.Context)

	resV, err := reqFn.Call(ctxV, reqV)
	res := resV.Interface()
	if err != nil {
		return res, err
	}

	if opts.Target != nil {
		err := copier.Copy(opts.Target, res)
		if err != nil {
			log.Print(err)
		}
	}

	return res, nil
}

func (e *Engine) Broadcast(notification interface{}, options ...decouple.CallOption) ([]interface{}, error) {
	notV := reflect.ValueOf(notification)
	notT, err := types.ParseType(notV.Type())
	if err != nil {
		panic(err)
	}
	subFn := e.Container.GetSubscribeFunctions(notT.ValueType)
	if len(subFn) == 0 {
		return nil, nil
	}
	opts := decouple.NewCallOptions(options)

	var allErrs error
	var allRets []interface{}
	ctxV := reflect.ValueOf(opts.Context)
	for _, fn := range subFn {
		val, err := fn.Call(ctxV, notV)
		if val.IsValid() && val.Kind() == reflect.Struct {
			allRets = append(allRets, val.Interface())
		}
		if err != nil {
			allErrs = multierr.Append(allErrs, err)
		}
	}

	return allRets, allErrs
}

var _ decouple.Call = (*Engine)(nil)

func NewEngine(container *decouple.Container) *Engine {
	return &Engine{
		Container: container,
	}
}
