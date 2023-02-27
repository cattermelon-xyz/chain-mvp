package machines

type VetoVote struct{}

func (this *VetoVote) Desc() string {
	return "This type of voting contains only 1 option is Veto. If the VETO voting result is smaller than the threshold value, this proposal will be passed"
}
