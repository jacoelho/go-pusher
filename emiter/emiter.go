package emiter

// An Emiter delievers events to a channel
type Emiter interface {
	Events() <-chan []byte
}
