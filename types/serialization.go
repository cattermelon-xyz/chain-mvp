package types

import "github.com/hectagon-finance/chain-mvp/types/event"

/*
* All neccessary information to reload a Mission
 */
type MissionData struct {
	Id       string
	Title    string
	Fulltext string
	Owner    string

	StartChkPId   string
	CurrentChkPId string
	IsStarted     bool
	IsActive      bool
}

/**
* Resurrect the data back to memory
 */
func (this *MissionData) Unmarshal() *Mission {
	return &Mission{}
}

/*
* All neccessary information to reload a CheckPoint
 */
type CheckPointData struct {
	Id               string
	Title            string
	Description      string
	FallbackId       uint64
	ChildrenId       []string
	LastBlockToVote  uint64
	LastBlockToTally uint64
	OutputEventData  event.EventData
	VoteMachineType  string
	VoteMachine      []byte
}

/**
* All neccessary information to reload the memory
 */
type SnapshotData struct {
	MissionDatas    []MissionData
	CheckPointDatas []CheckPointData
}

/**
* Rebuild the memory state from a
* Snapshot{AtBlock uint64, AtBlockHash uint64, Data []SnapshotData}
 */
func buildMemStateFromSnapshot(s Snapshot) {
	currentMemmoryAtBlocNo = s.AtBlock
	currentMemoryAtBlockHash = s.AtBlockHash
	Missions = make([]*Mission, 0)
	CheckPoints = make([]*CheckPoint, 0)
	for _, missionData := range s.Data.MissionDatas {
		m := missionData.unmarshal()
		Missions = append(Missions, m)
		for _, checkPointData := range s.Data.CheckPointDatas {
			if checkPointData.Id == string(m.id) {
				chkp := m.unmarshalFromCheckPointData(checkPointData)
				CheckPoints = append(CheckPoints, chkp)
			}
		}
	}
	// encode
	// var missions []Mission
	// decodeBuf := bytes.NewBuffer(s.Data)
	// decoder := gob.NewDecoder(decodeBuf)
	// err := decoder.Decode(&missions)
	// if err != nil {
	// 	log.Fatalf("Cannot decode data at block %d (%s)\n", s.AtBlock, s.AtBlockHash)
	// }
}

/*
* Periodically create a snapshot to provide other validator
 */
func createSnapshot() Snapshot {
	snapshotMissions := make([]MissionData, 0)
	snapshotCheckPoints := make([]CheckPointData, 0)
	for _, mission := range Missions {
		if mission.isActive == true {
			missionData := mission.marshal()
			snapshotMissions = append(snapshotMissions, missionData)
			for _, chkP := range CheckPoints {
				chkPData := chkP.marshal()
				snapshotCheckPoints = append(snapshotCheckPoints, chkPData)
			}
		}
	}
	return Snapshot{
		Data: SnapshotData{
			MissionDatas:    snapshotMissions,
			CheckPointDatas: snapshotCheckPoints,
		},
		AtBlock:     currentMemmoryAtBlocNo,
		AtBlockHash: currentMemoryAtBlockHash,
	}
}
