# `MVP` for the Chain

This project is the MVP for the new chain. These code focus on the application logic and try to deliver a `test-net` by the end of `May 2023`. This project do not focus on the production blockchain. A high level docs can be found [here](https://hectagondao.notion.site/D-Chain-Design-a7f071f3e7514191be453852a5675699)
<hr>

## Milestones:

1. Through `main.go`, you can express the workflow of building, managing and using of an `Initiative`. *Expected: end of Feb 2023*
2. Using a web interface, user can interactive with a go lang `web-server` to build, manage and use an `Initiative`. `Note:` using same tech stack with `cosmos` module to make the transition between phases easier. *Expected: end of Mar 2023*
3. Decentralize the `web-server` with `cosmos-sdk`. *Expected: end of May 2023*

## Documents

1. Q&A for chain design decision: [Link](assets/DesignQA.md)
2. Code convention: [Link](assets/Code.md)
3. Logic explaination: [Link](assets/Logic.md)

## TODO

- [x] Split cmd/demo into multiple programs
- [x] Change code structure to fit cosmos module coding recommendation
- [x] Build document structure
- [x] Change `Decision` to `Initiative` to match with business docs
- [ ] Implement the `CheckPoint` and `VotingMachine`
- [ ] The chain should emit an `Event` then some process will listen to it and do some work (aka `Enforcer`)
- [ ] Implement Observer pattern, NOTE: in the future, this will be replaced with Cosmos events
- [ ] Turn this project into a webserver
- [ ] Build web UI for initiative designer
- [ ] Build web UI for voter