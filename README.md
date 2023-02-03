# chain-mvp

MVP for the new chain.

This demo focus on the application logic, not yet the blockchain.
All docs could be found [here](https://hectagondao.notion.site/D-Chain-Design-a7f071f3e7514191be453852a5675699)

Code standard: [Link](https://github.com/golang-standards/project-layout)

# Convention

- File that contain `func main()` will be named `main.go`

# TODO

- [x] Split cmd/demo into multiple programs
- [ ] Read and apply cosmos module convention
- [ ] Add interfaces for pkg

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