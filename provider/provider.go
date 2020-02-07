package provider

import (
	"github.com/microapis/messages-core/channel"
	"github.com/microapis/events-api"
)

// Provider ...
type Provider interface {
	Approve(*events.Message) error
	Deliver(m *events.Message) error
	ToProvider() *channel.Provider
}
