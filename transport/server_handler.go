package transport

type ServerHandler interface {
	OnData([]byte)
}
