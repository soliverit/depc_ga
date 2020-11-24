package ga

import "strconv"

/*
	A record of the GA datasets state
 */
type GAState struct{
	//State of the data entries of the GA
	entityStates	[]GAStateRecord
	score			float32
	points			float32
	cost			float32
	scored			bool
}
/*
	Create a new GAState
 */
func CreateGAState(states []GAStateRecord)GAState{
	var gaState GAState 	= GAState{}
	/*
		Set properties from inputs
	 */
	gaState.entityStates	= states

	/*
		Defaults
	 */
	gaState.scored			= false
	/*
		Return stuff
	 */
	return gaState
}
/*
	Retrieve the cost
 */
func(gaState *GAState)Score()float32{
	return gaState.score
}
func(gaState *GAState)Scored()bool{
	return gaState.scored
}
/*
	Set the cost cache if it isn't set already
 */
func(gaState *GAState)SetScore(cost float32, points float32){
	if ! gaState.scored	{
		gaState.score 	= cost / points
		gaState.points	= points
		gaState.cost	= cost
		gaState.scored	= true
	}
}
func(gaState *GAState)RowState(idx int)*GAStateRecord{
	return &gaState.entityStates[idx]
}
/*
	Create a comma-delimited line from the score, points and cost
	unless these haven't been populated yet then return blank I suppose
 */
func(gaState *GAState)ToCSV()string{
	if gaState.scored{
		return 	strconv.FormatFloat(float64(gaState.score), 'f',4,32) + "," +
				strconv.FormatFloat(float64(gaState.points), 'f',4,32) + "," +
				strconv.FormatFloat(float64(gaState.cost), 'f',4,32)
	}
	return ""
}

