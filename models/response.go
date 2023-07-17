package models

type Response struct {
	Success bool
	Message string
	Data    interface{}
}
