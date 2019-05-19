package game

import (
	"sync"
)

type World struct {
	mu sync.Mutex

	entities []Entity
}

func (w *World) AddEntity(e Entity) {
	w.mu.Lock()
	w.entities = append(w.entities, e)
	w.mu.Unlock()
}

func (w *World) RemoveEntity(e Entity) {
	w.mu.Lock()
	for i, v := range w.entities {
		if v == e {
			w.entities = append(w.entities[:i], w.entities[i+1:]...)
		}
	}
	// w.entities = append(w.entities, entities)
	w.mu.Unlock()
}

type Entity interface {
	Tick(i int)
	// OnJoinWorld(w *World)
	// OnLeaveWorld(w *World)
	// OnSwitchWorld(oldW, newW *World)
}
