package ginject

import (
	"errors"
	"reflect"
)

type (

	// Injector is the actual interface for the injection functionality
	Injector interface {
		// Resolve transforms a lazy dep into a real one
		Resolve(name string)
		// RegisterByName enables the injector to get a service
		RegisterByName(name string, obj interface{})
		// Register enables the injector to get a service
		Register(service Service)
		// RegisterLazy enables the injector to lazyly register services
		RegisterLazy(name string, cb func() interface{})
		// Get returns previously registered services
		Get(name string, objPtr interface{}) error
		// Apply applies registered services on structs
		Apply(ptr interface{}) error
	}

	// Service is the interface for registering services with the injector
	Service interface {
		GetName() string
		Get() interface{}
	}

	srv struct {
		name string
		obj  interface{}
	}

	// Inj is the actual implementation of Injector
	Inj struct {
		deps map[string]reflect.Value
		lazy map[string]func() interface{}
	}
)

// ################### Service

// IService creates a service obj
func IService(name string, obj interface{}) Service {
	return srv{name: name, obj: obj}
}

// GetName returns the name of the service
func (srv srv) GetName() string {
	return srv.name
}

// Get returns the actual service
func (srv srv) Get() interface{} {
	return srv.obj
}

// #################### Injector

// New creates a new injector
func New() Injector {

	return &Inj{
		deps: make(map[string]reflect.Value),
		lazy: make(map[string]func() interface{}),
	}
}

// RegisterByName enables the injector to get a service
func (i *Inj) RegisterByName(name string, obj interface{}) {
	i.deps[name] = inflect(obj)
}

// Register enables the injector to get a service
func (i *Inj) Register(service Service) {
	i.RegisterByName(service.GetName(), service.Get())
}

// RegisterLazy enables the injector to get a function returnin a service on first use
func (i *Inj) RegisterLazy(name string, cb func() interface{}) {
	i.lazy[name] = cb
}

// Get returns previously registered services
func (i *Inj) Get(name string, objPtr interface{}) error {

	objType := reflect.TypeOf(objPtr)

	if objType.Kind() != reflect.Ptr {
		return errors.New("given obj is not a pointer(1) -> " + objType.Kind().String())
	}

	i.resolveLazy(name)
	val, ok := i.deps[name]
	if !ok {
		return errors.New("dependency " + name + " not present")
	}

	if objType == val.Type() {

		vof := reflect.ValueOf(objPtr)
		if !vof.CanSet() {
			vof = vof.Elem()
		}

		if !vof.CanSet() {
			return errors.New("could not set to " + vof.Kind().String() + " | " + vof.Type().String())
		}

		vof.Set(val.Elem())
		return nil
	}

	if objType.Kind() == reflect.Ptr {
		objType = objType.Elem()
	}

	if objType.Kind() == reflect.Interface && val.Type().Implements(objType) {

		vof := reflect.ValueOf(objPtr).Elem()
		vof.Set(val)
		return nil
	}

	return errors.New("could not apply illegal type on given type")
}

// Apply applies registered services on structs
func (i *Inj) Apply(ptr interface{}) error {

	v := reflect.ValueOf(ptr)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return errors.New("can only appy on structs")
	}

	t := v.Type()

	for j := 0; j < v.NumField(); j++ {

		f := v.Field(j)
		sf := t.Field(j)
		if f.CanSet() {

			theTag := sf.Tag.Get("inject")
			i.resolveLazy(theTag)

			val, ok := i.deps[theTag]
			if !ok {
				return errors.New("dependency " + theTag + " was unknown")
			}

			ft := f.Type()

			if !val.IsValid() {
				return errors.New("value can not be applied to " + ft.String())
			}

			if ft == val.Type() {

				f.Set(val)
				continue
			}

			if ft.Kind() == reflect.Interface && val.Type().Implements(ft) {

				f.Set(val)
				continue
			}

			return errors.New("type missmatch! Cannot apply " + val.Type().String() + " to " + ft.String())
		}

	}

	return nil
}

// Resolve transforms a lazy dep into a real one
func (i *Inj) Resolve(name string) {
	i.resolveLazy(name)
}

// ################### Helpers

func (i *Inj) resolveLazy(name string) {

	_, ok := i.deps[name]
	if !ok {

		fnc, ok := i.lazy[name]
		if ok {

			i.deps[name] = inflect(fnc())
			delete(i.lazy, name)
		}
	}
}

func inflect(obj interface{}) reflect.Value {

	return reflect.ValueOf(obj)
}
