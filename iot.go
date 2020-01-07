package messagesiot

const (
	// Channel ...
	Channel = "iot"
)

// Message ...
type Message struct {
	Topic    string            `json:"topic"`
	Payload  map[string]string `json:"payload"`
	Provider string            `json:"provider"`
}
