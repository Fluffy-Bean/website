package broker

import (
	"reflect"
)

type Handler any

type Broker struct {
	handlers map[string][]reflect.Value
	events   []any
}

func NewBroker() *Broker {
	return &Broker{
		handlers: make(map[string][]reflect.Value),
	}
}

func (e *Broker) RegisterHandler(handler Handler) {
	typeofHandler := reflect.TypeOf(handler)

	if typeofHandler.NumIn() != 1 {
		panic("handler must have one input parameter")
	}
	if typeofHandler.NumOut() != 0 {
		panic("handler must have no output parameters")
	}

	name := eventTypeOfName(typeofHandler.In(0))
	e.handlers[name] = append(e.handlers[name], reflect.ValueOf(handler))
}

func (e *Broker) BroadcastEvent(event any) {
	typeofEvent := reflect.TypeOf(event)

	name := eventTypeOfName(typeofEvent)

	handlers, ok := e.handlers[name]
	if !ok {
		return
	}

	args := []reflect.Value{
		reflect.ValueOf(event),
	}

	for _, handler := range handlers {
		handler.Call(args)
	}
}

func eventTypeOfName(typeOf reflect.Type) string {
	var name string

	if typeOf.Kind() == reflect.Ptr {
		name += "*"
		typeOf = typeOf.Elem()
	}

	return name + typeOf.PkgPath() + "." + typeOf.Name()
}
