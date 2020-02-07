package events

const (
	// Channel ...
	Channel = "events"
)

// Message ...
type Message struct {
	Topic    string            `json:"topic"`
	Payload  map[string]string `json:"payload"`
	Provider string            `json:"provider"`
	Status   string            `json:"status"`
}
