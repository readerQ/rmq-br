package rabbit

type MessageWriter interface {
	WriteMessage(msg Message) error
}
