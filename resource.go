package gorest

type Resource interface{
	New() Resource
	NewList() interface{}

	Id(id string)
	Valid() bool
}
