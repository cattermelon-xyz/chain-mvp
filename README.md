# `MVP` for the Chain

This project is the MVP for the new chain. These code focus on the application logic and try to deliver a `test-net` by the end of `May 2023`. This project do not focus on the production blockchain. A high level docs can be found [here](https://hectagondao.notion.site/D-Chain-Design-a7f071f3e7514191be453852a5675699)
<hr>

## Milestones

1. Through `main.go`, you can express the workflow of building, managing and using of an `Mission`. *Expected: end of Feb 2023*
2. Using a web interface, user can interactive with a go lang `web-server` to build, manage and use an `Mission`. `Note:` using same tech stack with `cosmos` module to make the transition between phases easier. *Expected: end of Mar 2023*

<ol>
  <li>Integrate Login & Logout web2 style</li>
  <li>Integrate wallet (connect wallet and use token to vote) & trigger Web2 enforcer actions (post twitter, interact with Airtable, discord/slack, send email...).</li>
</ol>

3. Decentralize the `web-server` with `cosmos-sdk`. *Expected: end of May 2023*

## Documents

1. Q&A for chain design decision: [Link](assets/DesignQA.md)
2. Code convention: [Link](assets/Code.md)
3. Logic explaination: [Link](assets/Logic.md)

## TODO

- [x] Split cmd/demo into multiple programs
- [x] Change code structure to fit cosmos module coding recommendation
- [x] Build document structure
- [x] Change `Decision` to `Mission` to match with business docs
- [x] Implement the `CheckPoint` and `VotingMachine`
- [x] Implement Cobra to build `htg` and `htg-client`
- [x] When a Tree travel to a new node, should call a hook with arguments as results from the last node voting. Possibly change VotingMachine{} and Node{} to match changes
- [x] Implement a mock API server
- [x] Implement a mock WebSocket server
- [x] Implement a mock Concensus goroutine to produce new Block
- [x] Integrate all mock goroutines
- [x] The chain should emit `Event`, this should reflect in the VotingMachine{} and be able to connect to outside with WebSocket
- [x] Support hiding voting result till revealation
- [x] Use Log instead of fmt.Println
- [x] Write testcase
- [x] Vote should not be recorded if fallbackAttempt is TRUE
- [x] Tally & Fallback should be called when a new block is produced
- [x] Vote should specify what CheckPoint it is voting, to prevent unintended vote
- [ ] Fix Mission, CheckPoint usage of pointer to enable marshal and unmarshal
- [ ] Turn the program into Blockchain friendly format (command extraction, sync, build memory from Snapshot and build memory from blockdata)
- [ ] Implement Polling, note that UpVote is just a special case of Polling
- [ ] Turn program into logic(mem, input) with each `logic` is assigned an `id` that the Event would emit
- [ ] Implement client commandline
- [ ] Demo ready
- [ ] Implement VetoVote
- [ ] Implement RankChoiceVote
- [ ] Implement SingleChoiceVote
- [ ] Define `proto` strategy how to use proto efficiently without absusing and lost project compartment chracteristic
- [ ] Implement `proto` to replace `struct` definition and json file
- [ ] Implement Observer pattern (this is a mechanism for `Enforcer` to hook up with `Event`) with a caution on memory leak & bloats on chain data, NOTE: in the future, this will be replaced with Cosmos events
- [ ] Implement server api call to match with all client call
- [ ] Integrate Web UI
- [ ] Website support MetaMask and Phantom to use ETH, BSC, Solana address to interact with the chain
