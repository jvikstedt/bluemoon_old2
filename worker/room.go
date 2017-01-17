package worker

type Entity interface {
	Update(delta float64) (bool, error)
	ID() int
	X() int
	Y() int
}

type Room interface {
	AddEvent(e Event)
}

type Event interface {
	Execute(room Room) error
}
