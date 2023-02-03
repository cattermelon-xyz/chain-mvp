# chain-mvp

MVP for the new chain.

This demo focus on the application logic, not yet the blockchain.
All docs could be found [here](https://hectagondao.notion.site/D-Chain-Design-a7f071f3e7514191be453852a5675699)

Code standard: [Link](https://github.com/golang-standards/project-layout)

# Convention

- File that contain `func main()` will be named `main.go`

# TODO

- [x] Split cmd/demo into multiple programs
- [x] Change code structure to fit cosmos module coding recommendation
- [ ] Implement Observer pattern, NOTE: in the future, this will be replaced with Cosmos events

# Example of voting
- Weight voting
- Veto voting

# Design note

- Tree contain Node
- Tree travel between Node, once a Node is active it will:
 - active its enforcer
 - active its checkpoint

 NOTE: change Enforcer to Event; the chain should emit an event then some process will listen to it and do some work
- Should Enforcer be written in the chain config or should it be something else written outside of the chain?

- Project structure
```
 - third_party // code copy from someone else placed here
 - client
  - cli // exposed command line
 - types // all types and internal logic goes here
 - rules // all voting rules defined here
```