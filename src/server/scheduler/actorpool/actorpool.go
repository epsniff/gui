package actorpool

import (
	"sync"

	"github.com/epsniff/gui/src/server/scheduler/peerstate"
	"github.com/lytics/grid"
)

// New peer queue.
func New(rebalance bool, ps peerstate.PeersState) *ActorPool {
	return &ActorPool{
		required:       map[string]*grid.ActorStart{},
		actorPoolState: newActorPoolState(ps),
		rebalance:      rebalance,
	}
}

type ActorPool struct {
	mu        sync.RWMutex
	required  map[string]*grid.ActorStart
	actorType string
	rebalance bool

	actorPoolState *ActorPoolState
}

// IsRequired actor.
func (ap *ActorPool) IsRequired(actor string) bool {
	ap.mu.RLock()
	defer ap.mu.RUnlock()

	_, ok := ap.required[actor]
	return ok
}

// ActorType of this peer queue.
func (ap *ActorPool) ActorType() string {
	ap.mu.RLock()
	defer ap.mu.RUnlock()

	return ap.actorType
}

// SetRequired flag on actor. If it's type does not match
// the type of previously set actors an error is returned.
func (ap *ActorPool) SetRequired(def *grid.ActorStart) error {
	ap.mu.Lock()
	defer ap.mu.Unlock()

	if !isValidName(def.Name) {
		return ErrInvalidActorName
	}

	if ap.actorType == "" {
		ap.actorType = def.Type
	}
	if ap.actorType != def.Type {
		return ErrActorTypeMismatch
	}
	ap.required[def.Name] = def
	return nil
}

// UnsetRequired flag on actor.
func (ap *ActorPool) UnsetRequired(actor string) {
	ap.mu.Lock()
	defer ap.mu.Unlock()

	delete(ap.required, actor)
	if len(ap.required) == 0 {
		ap.actorType = ""
	}
}

// Missing actors that are required but not registered.
func (ap *ActorPool) Missing() []*grid.ActorStart {
	ap.mu.RLock()
	defer ap.mu.RUnlock()

	var missing []*grid.ActorStart
	for name, def := range ap.required {
		if isReg := ap.actorPoolState.IsRegistered(name); !isReg {
			missing = append(missing, def)
		}
	}
	return missing
}

//NumRegistered returns the number of actors registered in this pool
func (ap *ActorPool) NumRegistered() int {
	return ap.actorPoolState.NumRegistered()
}

//NumRegisteredOn returns the number of actors registered in a peer in the pool
func (ap *ActorPool) NumRegisteredOn(peer string) int {
	return ap.actorPoolState.NumRegisteredOn(peer)
}

//NumRegistered returns the number of actors registered in this pool
func (ap *ActorPool) NumOptimisticallyRegistered() int {
	return ap.actorPoolState.NumOptimisticallyRegistered()
}

//NumRegisteredOn returns the number of actors registered in a peer in the pool
func (ap *ActorPool) NumOptimisticallyRegisteredOn(peer string) int {
	return ap.actorPoolState.NumOptimisticallyRegisteredOn(peer)
}

//IsRegistered returns if the actorName as been registered already.
func (ap *ActorPool) IsRegistered(actorName string) bool {
	return ap.actorPoolState.IsRegistered(actorName)
}

func (ap *ActorPool) IsOptimisticallyRegistered(actorName string) bool {
	return ap.actorPoolState.IsOptimisticallyRegister(actorName)
}

// Register the actor to the peer.
func (ap *ActorPool) Register(actor, peer string) error {
	if !isValidName(actor) {
		return ErrInvalidActorName
	}
	if !isValidName(peer) {
		return ErrInvalidPeerName
	}
	ap.actorPoolState.register(actor, peer)
	return nil
}

// OptimisticallyRegister an actor, ie: no confirmation has
// arrived that the actor is actually running on the peer,
// but it has been requested to run on the peer.
func (ap *ActorPool) OptimisticallyRegister(actor, peer string) error {
	if !isValidName(actor) {
		return ErrInvalidActorName
	}
	if !isValidName(peer) {
		return ErrInvalidPeerName
	}
	ap.actorPoolState.optimisticallyRegister(actor, peer)
	return nil
}

// Unregister the actor from its current peer.
func (ap *ActorPool) Unregister(actor string) error {
	if !isValidName(actor) {
		return ErrInvalidActorName
	}
	ap.actorPoolState.unregister(actor)
	return nil
}

// OptimisticallyUnregister the actor, ie: no confirmation has
// arrived that the actor is NOT running on the peer, but
// perhaps because of a failed request to the peer to start
// the actor it is known that likely the actor is not running.
func (ap *ActorPool) OptimisticallyUnregister(actor string) error {
	if !isValidName(actor) {
		return ErrInvalidActorName
	}
	ap.actorPoolState.optimisticallyUnregister(actor)
	return nil
}

//Status returns a struct that represents all the peer queue's internal states used
//for logging and debugging
func (ap *ActorPool) Status() *PeersStatus {
	return ap.actorPoolState.Status()
}

func isValidName(name string) bool {
	return name != ""
}
