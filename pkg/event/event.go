package event

var emitted []Event

type Event struct {
	Name string
	Args []string
}

func Emit(e Event) {
	emitted = append(emitted, e)
}
