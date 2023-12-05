package ga

type GAStateRecord struct {
	optionID int
}

func CreateGAStateRecord(optionID int) GAStateRecord {
	var stateRecord GAStateRecord
	stateRecord.optionID = optionID
	return stateRecord
}

/*
===

	Getters

===
*/
func (stateRecord *GAStateRecord) OptionID() int { return stateRecord.optionID }
