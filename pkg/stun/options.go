package stun

const (
	DEFAULT_LOGGING          = true
	DEFAULT_MAX_MSG_IN_QUEUE = 10
)

// Stun options struct
type StunOptions struct {
	logging bool
}

// Creates a new stun options
func NewStunOptions(logging bool) StunOptions {
	return StunOptions{logging: logging}
}

// Creates a new default stun options
func DefaultStunOptions() StunOptions {
	return NewStunOptions(DEFAULT_LOGGING)
}

// Client stun options struct
type ClientStunOptions struct {
	logging       bool
	maxMsgInQueue int
}

// Creates a new client stun options
func NewClientStunOptions(logging bool, maxMsgInQueue int) ClientStunOptions {
	return ClientStunOptions{logging: logging, maxMsgInQueue: maxMsgInQueue}
}

// Creates a new default client stun options
func DefaultClientStunOptions() ClientStunOptions {
	return NewClientStunOptions(DEFAULT_LOGGING, DEFAULT_MAX_MSG_IN_QUEUE)
}
