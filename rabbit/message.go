package rabbit

type Message struct {
	Queue string
	Index int
	Body  []byte
}
