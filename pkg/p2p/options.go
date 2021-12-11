package p2p

const (
	DEFAULT_MAX_MSG_IN_QUEUE = 10
	DEFAULT_SECONDS_TIMEOUT  = 10
)

// Peer options struct
type PeerOptions struct {
	maxMsgInQueue int
	timeout       int
}

// Creates a new peer options
func NewPeerOptions(maxMsgInQueue int, timeout int) PeerOptions {
	return PeerOptions{maxMsgInQueue: maxMsgInQueue, timeout: timeout}
}

// Creates a new default peer options
func DefaultPeerOptions() PeerOptions {
	return NewPeerOptions(DEFAULT_MAX_MSG_IN_QUEUE, DEFAULT_SECONDS_TIMEOUT)
}
