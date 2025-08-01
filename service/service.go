package service

type Handler interface {
	Read(msg string) (response string, finished bool)
}
