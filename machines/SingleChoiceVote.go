package machines

type SingleChoiceVote struct {
	Threshold int
}

func (this *SingleChoiceVote) Desc() string {
	return "User can choose 1 of 2 options. The option which has voting result is equivalent to or bigger than the threshold value will be passed."
}
