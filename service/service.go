package service

func PrepareLookupSet(s string) map[rune]struct{} {
	lookupSet := make(map[rune]struct{}, len(s))
	for _, r := range s {
		lookupSet[r] = struct{}{}
	}
	return lookupSet
}

// Handler defines the interface for a state in our WebSocket state machine.
type Handler interface {
	// Read processes a message and returns a response and the next handler.
	// If the state does not change, it should return itself.
	Read(msg string) (response string, next Handler)
}
