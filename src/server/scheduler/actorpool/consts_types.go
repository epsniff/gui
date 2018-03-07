package actorpool

import "errors"

var (
	ErrEmpty              = errors.New("empty")
	ErrInvalidActorName   = errors.New("invalid actor name")
	ErrInvalidPeerName    = errors.New("invalid peer name")
	ErrActorTypeMismatch  = errors.New("actor type mismatch")
	ErrActorNotRegistered = errors.New("actor not registered")
)
