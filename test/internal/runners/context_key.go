package runners

import "reflect"

type typeKey[V any] struct{}

func (typeKey[V]) isKey() {
	// Do nothing as it is a stub
}

func (typeKey[V]) Type() reflect.Type {
	var v V
	return reflect.TypeOf(v)
}

func NewKey[V any](_ ...V) StorageKey {
	return typeKey[V]{}
}
