package ga

/*
	A record of the GA datasets state
 */
type GAState struct{
	//State of the data entries of the GA
	entityStates	[]GAStateRecord
	score			float32
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
func(gaState *GAState)SetScore(cost float32){
	if ! gaState.scored	{
		gaState.score 	= cost
		gaState.scored	= true
	}
}
func(gaState *GAState)RowState(idx int)*GAStateRecord{
	return &gaState.entityStates[idx]
}

