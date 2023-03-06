package machines

import (
	"encoding/json"
	"fmt"

	"github.com/hectagon-finance/chain-mvp/types"
)

// TODO: A Polling always has 2 option: one is the next CheckPoint and another is the Fallback
type Polling struct {
	choices            []string
	choicesResult      map[uint64]uint64
	fallbackOptionId   uint64
	started            bool
	startedBlock       uint64
	noOfOptions        uint64
	tallyBeforeBlockNo uint64
	voteBeforeBlockNo  uint64
	maxChoicePerVote   uint64 // must input from NewPolling
	tallyResult        []byte
	lastTalliedBlock   uint64
}

/**
* Return new Polling. Require maxChoicePerVote > 0
 */
func NewPolling(predefinedChoices []string,
	maxChoicePerVote uint64, tallyBeforeBlockNo uint64, voteBeforeBlockNo uint64,
	noOfOptions uint64, fallBackOptionId uint64) *Polling {
	if maxChoicePerVote == 0 || tallyBeforeBlockNo == 0 {
		return nil
	}
	choices := []string{} // reset any pre-defined choices
	choicesResult := make(map[uint64]uint64)
	if predefinedChoices != nil && len(predefinedChoices) > 0 {
		for idx, c := range predefinedChoices {
			choices = append(choices, c)
			choicesResult[uint64(idx)] = 0
		}
	}
	return &Polling{
		choices:            choices,
		choicesResult:      choicesResult,
		noOfOptions:        noOfOptions,
		fallbackOptionId:   fallBackOptionId,
		maxChoicePerVote:   maxChoicePerVote,
		started:            false,
		startedBlock:       0,
		voteBeforeBlockNo:  voteBeforeBlockNo,
		tallyBeforeBlockNo: tallyBeforeBlockNo,
		lastTalliedBlock:   0,
	}
}

/**
* Return human-readable description of the VotingMachine
 */
func (this *Polling) Desc() string {
	var desc string
	defaultDesc := fmt.Sprintf(`Voting is not started yet. Users can choose %d.
	User can vote only once. The poll will tally everytime user vote.
	The voting results will include %d options.`, this.maxChoicePerVote, this.maxChoicePerVote)
	if this.started == false {
		desc = defaultDesc
	} else {
		choices := ""
		for idx, c := range this.choices {
			choices += fmt.Sprintf("%d. %s\n", idx, c)
		}
		desc = fmt.Sprintf(`Users can vote maximum %d option of the following:
		%s
		User can vote only once before block #%d. 
		Voting started at block #%d. The poll will tally everytime user vote and the last block to tally is block #%d.
		The voting results will include %d options. 
		There are %d possible outcome of this polling.
		If no result is conclude then fallback at child #%d of the Node will be triggered.\n`,
			this.maxChoicePerVote, choices, this.startedBlock, this.voteBeforeBlockNo, this.lastTalliedBlock, this.maxChoicePerVote, this.noOfOptions, this.fallbackOptionId)
	}
	return desc
}

/**
* Params: tallyResult []byte, noOfOptions uint64, startedBlock uint64, fallbackOption uint64
* - noOfOptions are number of children that the tree can travel after voting conclude. For Polling, this is different from the option that user can choose
* - fallbackOption is the index of children that tree can travel if no result is given.
* Return `started` succeed. After this, the machine is ready for vote.
 */
func (this *Polling) Start(tallyResult []byte, noOfOptions uint64, startedBlock uint64, fallbackOption uint64) bool {
	this.started = true
	var err error
	if tallyResult != nil {
		var userChoices []string
		err = json.Unmarshal(tallyResult, &userChoices)
		if err != nil && len(userChoices) > 0 {
			if len(userChoices) <= int(this.maxChoicePerVote) {
				this.started = false
			} else {
				this.choices = []string{} // reset any pre-defined choices
				this.choicesResult = make(map[uint64]uint64)
				for idx, c := range userChoices {
					this.choices = append(this.choices, c)
					this.choicesResult[uint64(idx)] = 0
				}
			}
		} else {
			this.started = false
		}
	}
	if fallbackOption >= noOfOptions {
		this.started = false
	}
	this.noOfOptions = noOfOptions
	this.startedBlock = startedBlock
	return this.started
}

// Return the Readiness of the VotingMachine
func (this *Polling) IsStarted() bool {
	return this.started
}

/**
* Params: rawChoicesIdx []byte
* rawChoicesIdx must be []string
* Return True if the option is validated
 */
func (this *Polling) ValidateVote(rawChoicesIdx []byte) bool {
	currentBlockNumber := types.GetCurrentBlockNumber()
	if currentBlockNumber > this.voteBeforeBlockNo {
		return false
	}
	var choicesIdx []int
	err := json.Unmarshal(rawChoicesIdx, &choicesIdx)
	for _, val := range choicesIdx {
		if val >= len(this.choices) {
			return false
		}
	}
	if err != nil {
		return false
	}
	return true
}

/**
* Params: who string, rawChoicesIdx []byte
* rawChoicesIdx must be []int
* // TODO: what if we want to hide the who-vote-what from Validator?
* Return True if the option is validated
 */
func (this *Polling) Record(who string, rawChoicesIdx []byte) bool {
	if this.started == false {
		return false
	}
	var choicesIdx []int
	err := json.Unmarshal(rawChoicesIdx, &choicesIdx)
	if err != nil {
		return false
	}
	for _, val := range choicesIdx {
		this.choicesResult[uint64(val)] += this.VotingPower(who, rawChoicesIdx)
	}
	return true
}

/**
* Params: who string, rawChoicesIdx []byte
* rawChoicesIdx must be []int
* TODO: what if we want to hide the who-vote-what from Validator?
* Return the VotingPower of the vote
 */
func (this *Polling) VotingPower(who string, rawChoicesIdx []byte) uint64 {
	// TODO: connect to oracle to calculate voting power
	return 1
}

/**
* TODO: what if we want to hide the who-vote-what from Validator?
* TODO: Tally at the voting
* Function is called everytime someone vote
* Return true | false if Polling Should Tally
 */
func (this *Polling) ShouldTally() bool {
	if this.started == false {
		return false
	}
	currentBlock := types.GetCurrentBlockNumber()
	if this.tallyBeforeBlockNo > currentBlock {
		return false
	}
	return true
}

/**
* Params: who string, rawChoicesIdx []byte
* Return the Cost for the Vote
 */
func (this *Polling) Cost(who string, rawChoicesIdx []byte) uint64 {
	return 0
}

/**
* Params: blockNumber uint64
* TODO: what if we want to hide the who-vote-what from Validator?
* Return if Tallied is succeed.
 */
func (this *Polling) Tally() bool {
	if this.started == false {
		return false
	}
	currentBlock := types.GetCurrentBlockNumber()
	if currentBlock > this.tallyBeforeBlockNo {
		return false
	}
	// order the choicesResult with max on top
	// find top maxChoicePerVote of choicesResult
	this.lastTalliedBlock = currentBlock
	return true
}

/**
* Return the last successful Tallied Block
 */
func (this *Polling) GetLastTalliedBlock() uint64 {
	return this.lastTalliedBlock
}

// Return the Tally result, return nil []byte and NoOptionMade code if no option made.
func (this *Polling) GetTallyResult() ([]byte, uint64) {

	if this.started == false {
		return nil, types.NoOptionMade
	}
	if this.lastTalliedBlock != 0 {
	}
	return nil, 0
}
