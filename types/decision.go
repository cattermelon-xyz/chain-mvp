package types

type Decision struct {
	Title    string
	Fulltext string
	Start    *CheckPoint
	Id       Address
	Owner    Address
	State    *State
}

type State struct {
	IsStarted bool
	Current   *CheckPoint
}

type iDecision interface {
	create(Decision) string
}

func create(title string, fulltext string, start *CheckPoint) *Decision {
	return nil
}

func (this *Decision) edit(d Decision) bool {
	return false
}

func (this *Decision) delete(d Decision) bool {
	return false
}

func (this *Decision) start(d Decision) bool {
	return false
}

func (this *Decision) stop(d Decision) bool {
	return false
}

func (this *Decision) pause(d Decision) bool {
	return false
}

func (this *Decision) resume(d Decision) bool {
	return false
}
