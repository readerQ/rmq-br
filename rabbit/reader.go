package rabbit

type MessageReader interface {
	ReadMessage() (Message, bool, error)
}
