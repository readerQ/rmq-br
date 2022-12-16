package rabbit

type Message struct {
	Queue       string
	ContentType string
	Index       int
	Body        []byte
}
