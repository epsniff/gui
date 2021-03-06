package leader

import (
	"context"
	"time"

	"github.com/epsniff/gui/src/lib/logging"
	"github.com/lytics/grid"
)

const defautTimeout = 2 * time.Second

// LeaderActor is the scheduler to create and watch
// the workers but the work comes from http requests
type LeaderActor struct {
	client *grid.Client
	cfg    *Cfg
}

func New(client *grid.Client, cfg *Cfg) *LeaderActor {
	if cfg.timeout <= time.Duration(0) {
		cfg.timeout = defautTimeout
	}
	return &LeaderActor{
		client: client,
		cfg:    cfg,
	}
}

// Act checks for peers, ie: other processes running this code,
// in the same namespace and start the worker actor on one of them.
func (a *LeaderActor) Act(c context.Context) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	logging.Logger.Infof("Starting Leader Actor")
	defer logging.Logger.Infof("Exiting Leader Actor")

	existing := make(map[string]bool)
	for {
		select {
		case <-c.Done():
			return
		case <-ticker.C:
			// Ask for current peers.
			peers, err := a.client.Query(a.cfg.timeout, grid.Peers)
			if err != nil {
				logging.Logger.Errorf("got an error quering peers: err:%v", err)
				return
			}

			// Check for new peers.
			for _, peer := range peers {
				if existing[peer.Name()] {
					continue
				}

				// Define a worker.
				existing[peer.Name()] = true
				start := grid.NewActorStart("worker-%d", len(existing))
				start.Type = "worker"

				// On new peers start the worker.
				//TODO retry.X
				_, err := a.client.Request(a.cfg.timeout, peer.Name(), start)
				if err != nil {
					logging.Logger.Errorf("got an error starting actor:%v peer:%v err:%v", start, peer, err)
					return
				}
			}
		}
	}
}
