# Logic and UseCases

## Usecases

There are 2 main users of the app: `initiative designer` is one who design the all checkpoints, events and their enforcers. `voter` is the one who make decision hence trigger the `initiative` to function.

### Initiative Designer

- Create
- Update
- Delete
- Start
- Stop
- Pause
- Resume

### Voter

- Vote

## Design Logic

- An `Intitative` is a tree of `Node`; each node has an `Event` and/or a `Checkpoint`. `Event` will be emitted once the node is become the `Current` node of the tree.
- `CheckPoint` help `Intiative` travel between `Node`
- `VotingMachine`
  - A `Checkpoint` has a `VotingMachine` which has its own logic and store its own data. The `VotingMachine`.`record()` will be called upon a `Who`.`vote()`.
  - This is VotingMachine Interface

  ```go
    // Record to record the data; Tally to take action from the data; TallyAt return the timestamp to active Tally
    type VotingMachine interface {
      // Describe the rule of the vote
      Desc() string
      // Record the data: who <string> choose option <int>
      Record(string, int)
      // Calculate the voting power of the vote
      VotingPower(string, int) int
      // Cost of the Vote
      Cost(string, int) int
      // Tally the vote
      Tally()
      // When to tally the vote; if TallyAt() != -1 then it can only tally ONCE
      TallyAt() int
      // Return the Tally result, return NoTallyResult if no option is made
      GetTallyResult() int
    }
  ```
  
  - A `VotingMachine` defines when it should `Tally()` the vote result by its `TallyAt() timestamp`. If `TallyAt()` return `const TallyAfterVote` then `Tally()` will be executed after every vote, else it will execute once the timestamp is reached:

  ```go
    if Initiative.Current().IsTallied() == false && 
    Initiative.Current().TallyAt() >= Block.timestamp 
    { 
      Initiative.Current().Tally()
      Remove Initiative from the List 
    }
  ```

  - `VotingMachine.IsTallied() bool` return whether the VotingMachine already tally its voting result. This function ensure vote result can be re-tally if there is something wrong with the timestamp. NOTE: how Cosmos handle time and timestamp?
- An `Enforcer` listens to an `Event` and fired upon an event is emitted.

## Type of voting

- SingleChoiceVote
- Polling
- UpVote
- VetoVote
- RankVote
