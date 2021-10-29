package controller

type Context interface {
	Bind(interface{}) error
	Param(string) string
	JSON(int, interface{})
}
