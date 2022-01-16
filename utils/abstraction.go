package utils

type WsPusher interface {
	Push(interface{}) error
}
