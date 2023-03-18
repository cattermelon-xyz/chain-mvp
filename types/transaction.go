package types

import (
	"bytes"
	"encoding/gob"
)

type Command string

const (
	CMDCreateMission Command = "Create"
	CMDStartMission  Command = "Start"
	CMDPauseMission  Command = "Pause"
	CMDStopMission   Command = "Stop"
	CMDVote          Command = "Vote"
)

type CreateMissionCommandData struct {
	ChkPs       []CheckPoint
	StartChkPId string
	Title       string
	Fulltext    string
	Owner       Address
}

type StartMissionCommandData struct {
	MissionIdStr string
}

type PauseMissionCommandData struct {
	MissionIdStr string
	Description  string
}

type StopMissionCommandData struct {
	MissionIdStr string
	Description  string
}

type VoteCommandData struct {
	option []byte
}

type Transaction struct {
	Who  string
	Cmd  Command
	Data []byte
}

type Snapshot struct {
	Data        SnapshotData
	AtBlock     uint64
	AtBlockHash string
}

var currentMemmoryAtBlocNo uint64
var currentMemoryAtBlockHash string

/*
* TODO: Instead of 1 mem,
* We use: rawMissions []byte, rawCheckPoints []byte to save computing effort
* Instead of resurrected all Mission, we only touch the mentioned one
* Only use this function to sync memory state
 */
func sync(s Snapshot, blockDatas []BlockData) {
	buildMemStateFromSnapshot(s)           // memory is initialized properly with all global variables and so on
	for _, blockData := range blockDatas { // loop through all blockData
		buildMemStateFromBlockData(blockData)
	}
	// encode
}

/*
* TODO: Btw, if the validator is running for a while and everything is on the memory,
* should we resurrect anything at all?
 */
func buildMemStateFromBlockData(blockData BlockData) {
	// blockNumber := blockData.CurrentBlockNumber
	// blockHash := blockData.BlockHash
	transactions := blockData.Transactions
	for _, t := range transactions { // loop through all transaction & extract command
		switch t.Cmd {
		case CMDCreateMission:
			buff := bytes.NewBuffer(t.Data)
			createMissionData := CreateMissionCommandData{}
			enc := gob.NewDecoder(buff)
			err := enc.Decode(&createMissionData)
			if err != nil {

			}

			break
		case CMDPauseMission:
			break
		case CMDStartMission:
			break
		case CMDStopMission:
			break
		case CMDVote:
			break
		}
	}
}
