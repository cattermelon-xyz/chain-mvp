package machines

type Polling struct{}

func (this *Polling) Desc() string {
	return "Users can choose more than 1 option. The voting results can include 1 option or a list of options based on the voting condition setting."
}
