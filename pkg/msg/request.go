package msg

// Stun to Peer, Peer to Stun and
// Peer to Peer message request interface
type MsgRequest struct {
	Action   string `json:"action"`
	Peername string `json:"peername"`
	Message  string `json:"message"`
}

// Creates a new msg request
func NewMsgRequest(action string, peername string, message string) MsgRequest {
	return MsgRequest{
		Action:   action,
		Peername: peername,
		Message:  message,
	}
}
