package bluemoon

import "sync"

type IDGen struct {
	next  int
	nLock sync.Mutex
}

func NewIDGen() *IDGen {
	return &IDGen{
		next: 1,
	}
}

func (idg *IDGen) Next() int {
	idg.nLock.Lock()
	defer idg.nLock.Unlock()

	id := idg.next
	idg.next++
	return id
}
