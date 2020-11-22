package ga

import (
	"../csv"
)

type BaseGA struct{
	// State history
	states 			[]GAState
	//Random state seed
	seed			int
	//Number of steps
	iterations 		int
	//Input data
	data			*csv.BuildingReader
	//Scorer
	scorer			*IScorer
	//Hardness
	Hardness		float32
}


