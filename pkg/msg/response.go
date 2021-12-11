package msg

// Stun to Peer, Peer to Stun and
// Peer to Peer message response interface
type MsgResponse struct {
	Action   string `json:"action"`
	HasError bool   `json:"has_error"`
	Peername string `json:"peername"`
	Message  string `json:"message"`
}

// Creates a new msg response
func NewMsgResponse(action string, hasError bool, peername string, message string) MsgResponse {
	return MsgResponse{
		Action:   action,
		HasError: hasError,
		Peername: peername,
		Message:  message,
	}
}
