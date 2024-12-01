package status

type Class int16

const (
	ClassIncomplete Class = iota
	ClassCompleted
	ClassCancelled
)
