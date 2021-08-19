package pkg

import (
	"decouple/pkg/types"
	"fmt"
	"reflect"
	"sync"
)

type Container struct {
	NameByType map[reflect.Type]string
	TypeByName map[string]reflect.Type

	mNames sync.RWMutex

	RequestMap  map[string]*types.Function
	mRequestMap sync.RWMutex

	SubsMap  map[string][]*types.Function
	mSubsMap sync.RWMutex
}

var _ Decouple = (*Container)(nil)

func NewContainer() *Container {
	cnt := &Container{
		NameByType: map[reflect.Type]string{},
		TypeByName: map[string]reflect.Type{},
		RequestMap: map[string]*types.Function{},
		SubsMap:    map[string][]*types.Function{},
	}
	return cnt
}

func (c *Container) Request(fn interface{}) {
	fnReq, err := types.ParseFunction(fn)
	if err != nil {
		panic(err)
	}
	if err := fnReq.ValidateRequest(); err != nil {
		panic(err)
	}
	name := c.GetTypeName(fnReq.Argument.ValueType)

	c.mRequestMap.Lock()
	defer c.mRequestMap.Unlock()
	c.RequestMap[name] = fnReq
}

func (c *Container) Subscribe(fn interface{}) {
	sFn, err := types.ParseFunction(fn)
	if err != nil {
		panic(err)
	}
	if err := sFn.ValidateSubscribe(); err != nil {
		panic(err)
	}
	name := c.GetTypeName(sFn.Argument.ValueType)

	c.mSubsMap.Lock()
	defer c.mSubsMap.Unlock()

	list, _ := c.SubsMap[name]
	c.SubsMap[name] = append(list, sFn)
}

func (c *Container) SetProtoName(model interface{}, name string) {
	c.SetTypeName(reflect.TypeOf(model), name)
}

func (c *Container) GetProtoName(model interface{}) string {
	return c.GetTypeName(reflect.TypeOf(model))
}

func (c *Container) SetTypeName(rt reflect.Type, name string) {
	t, err := types.ParseTypeNamed(rt, name)
	if err != nil {
		panic(err)
	}

	c.mNames.RLock()
	et, ok := c.TypeByName[t.Name]
	c.mNames.RUnlock()

	if t.ValueType != et {
		panic(fmt.Errorf("type mismatch for name %s %s != %s", t.Name, t.ValueType.Name(), et.Name()))
	}

	if ok {
		return
	}

	c.mNames.Lock()
	defer c.mNames.Unlock()
	c.NameByType[t.ValueType] = t.Name
	c.TypeByName[t.Name] = t.ValueType
}

func (c *Container) GetTypeName(rt reflect.Type) string {
	t, err := types.ParseTypeNamed(rt, "")
	if err != nil {
		panic(err)
	}

	c.mNames.RLock()
	name, ok := c.NameByType[t.ValueType]
	c.mNames.RUnlock()

	if ok {
		return name
	}
	return t.Name
}

func (c *Container) GetRequestFunction(t reflect.Type) *types.Function {
	name := c.GetTypeName(t)

	c.mRequestMap.RLock()
	defer c.mRequestMap.RUnlock()
	rf, _ := c.RequestMap[name]

	return rf
}

func (c *Container) GetSubscribeFunctions(t reflect.Type) []*types.Function {
	name := c.GetTypeName(t)

	c.mSubsMap.RLock()
	defer c.mSubsMap.RUnlock()
	sl, _ := c.SubsMap[name]

	return sl
}
