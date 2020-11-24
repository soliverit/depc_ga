package ga

import (
	"../csv"
)
type EPCScorer struct{
	IScorer
	BaseScorer
	data 	*csv.BuildingReader
}
func CreateEPCScorer(data *csv.BuildingReader) *EPCScorer{
	var epcScorer EPCScorer = EPCScorer{}
	epcScorer.data			= data
	epcScorer.Description	= "RdSAP EPC Scorer"

	return &epcScorer
}
/*
	Score the input GAState
 */
func(epcScorer *EPCScorer) Score(gaState *GAState)float32{
	/*
		If the Score has been calculated already, return it.
	 */
	if gaState.scored{
		return gaState.score
	}
	var data *csv.BuildingReader = epcScorer.data
	var cost 		float32
	var points 		float32
	var building 	*csv.Building
	var stateRecord	*GAStateRecord
	for i := 0; i < data.Length(); i++{
		building 	= epcScorer.data.Building(i)
		stateRecord	= gaState.RowState(i)
		if stateRecord.EfficiencyIndex() != -1{
			cost 	+= building.Cell(stateRecord.CostIndex())
			points 	+= building.Cell(stateRecord.EfficiencyIndex())
		}
	}
	/*
		Cache the total cost in the GAState and return it.

		TODO: 	This should be a map with multiple costs. Eventually,
				change GAState.cost to map[string]float32
	 */
	gaState.SetScore(cost, points)
	return gaState.score
}


