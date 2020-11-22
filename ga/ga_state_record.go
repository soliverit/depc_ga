package ga
type GAStateRecord struct{
	efficiencyIndex int
	costIndex		int
}
func CreateGAStateRecord(efficiencyIndex int, costIndex int)GAStateRecord{
	var stateRecord GAStateRecord
	stateRecord.efficiencyIndex	= efficiencyIndex
	stateRecord.costIndex		= costIndex
	return stateRecord
}
/*===
	Getters
===*/
func(stateRecord *GAStateRecord)EfficiencyIndex()int{
	return stateRecord.efficiencyIndex
}
func(stateRecord *GAStateRecord)CostIndex()int{
	return stateRecord.costIndex
}