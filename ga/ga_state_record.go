package ga

import "strconv"

type GAStateRecord struct {
	efficiencyIndex int
	costIndex       int
}

func CreateGAStateRecord(efficiencyIndex int, costIndex int) GAStateRecord {
	var stateRecord GAStateRecord
	stateRecord.efficiencyIndex = efficiencyIndex
	stateRecord.costIndex = costIndex
	return stateRecord
}

/*===
	Getters
===*/
func (stateRecord *GAStateRecord) EfficiencyIndex() int {
	return stateRecord.efficiencyIndex
}
func (stateRecord *GAStateRecord) CostIndex() int {
	return stateRecord.costIndex
}
func (stateRecord *GAStateRecord) IsDefault() bool {
	return stateRecord.efficiencyIndex == -1
}
func (stateRecord *GAStateRecord) IsModification() bool {
	return stateRecord.efficiencyIndex != -1
}
func (stateRecord *GAStateRecord) ToString() string {
	return "Cost ID: " + strconv.Itoa(stateRecord.costIndex) + "\tEff ID: " + strconv.Itoa(stateRecord.efficiencyIndex)
}
