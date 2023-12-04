package ga

import (
	"../csv"
	"../helpers"
	"math/rand"
	"sort"
)

const LOG_PATH string = "c:/repos/ga_results.txt"

type EpcGA struct {
	BaseGA
	data   *csv.BuildingReader
	scorer *EPCScorer
	//Building state options
	buildingRetorfits [][]int
	effHeaders        []string
	effHeaderIndices  []int
	costHeaders       []string
	costHeaderIndices []int
	population        []GAState
	initialPopulation int
	CrossoverRate     float32
	childCount        int
	maxPopulation     int
}

/*
Create a new epcGAA, an extension of Baseg for residential EPC data
*/
func CreateEpcGA(bReader *csv.BuildingReader, iterations int, maxPopulation int, packages []string) *EpcGA {
	/*
		Create the output
	*/
	var epcGA EpcGA
	/*
		Set properties from inputs
	*/
	epcGA.data = bReader
	epcGA.iterations = iterations
	epcGA.scorer = CreateEPCScorer(epcGA.data)
	/*
		Do default internal stuff

		-------- Best ---
		crossoverRate:	0.15
		childCount:		10
		maxPopulation:	30
		hardness:		0.1
	*/

	epcGA.maxPopulation = maxPopulation
	epcGA.population = make([]GAState, epcGA.maxPopulation)
	epcGA.CrossoverRate = 0.15
	epcGA.childCount = 10
	epcGA.Hardness = 0.1 //Best 0.1 with crossoverRate of 0.15
	//Prepare arrays TODO: See below, put in Struct
	epcGA.costHeaders = make([]string, len(packages))
	epcGA.effHeaders = make([]string, len(packages))
	epcGA.effHeaderIndices = make([]int, len(packages))
	epcGA.costHeaderIndices = make([]int, len(packages))
	epcGA.population = make([]GAState, epcGA.maxPopulation)
	/*
		Generate label/column indices

		TODO: Make this into its own Struct
	*/
	for i := 0; i < len(packages); i++ {
		epcGA.effHeaders[i] = packages[i] + "-Eff"
		epcGA.costHeaders[i] = packages[i] + "-Cost"
		epcGA.effHeaderIndices[i] = epcGA.data.EasyReader.ColumnNameToIndex(epcGA.effHeaders[i])
		epcGA.costHeaderIndices[i] = epcGA.data.EasyReader.ColumnNameToIndex(epcGA.costHeaders[i])
	}

	/*
		Do Building state validity sets
	*/
	var buildingCount int
	var tempBuilding *csv.Building
	var tempRetrofitIDSet []int
	for buildingIDX := 0; buildingIDX < buildingCount; buildingIDX++ {
		/*
			Index all Building potential states
		*/
		tempRetrofitIDSet = make([]int, 0)
		tempBuilding = epcGA.data.Building(buildingIDX)
		for headerIDX := 0; headerIDX < len(epcGA.effHeaderIndices); headerIDX++ {
			if tempBuilding.Cell(epcGA.effHeaderIndices[headerIDX]) != -1 {
				tempRetrofitIDSet = append(tempRetrofitIDSet, epcGA.effHeaderIndices[headerIDX])
			}
		}

	}
	/*
		Send it home
	*/
	return &epcGA
}
func (g *EpcGA) Best() float32 {
	var building *csv.Building
	/*
		Get header positions
	*/
	var cost float32 = 0.0
	var score float32 = 0.0
	var tempScore float32 = 0.0
	var tempRatingIDX int = 0
	var testScore float32
	for i := 0; i < g.data.Length(); i++ {
		building = g.data.Building(i)
		for j := 0; j < len(g.effHeaders); j++ {
			testScore = building.Cell(g.effHeaderIndices[j]) / building.Cell(g.costHeaderIndices[j])
			if testScore > tempScore {
				tempScore = testScore
				tempRatingIDX = j
			}
		}
		score += building.Cell(g.costHeaderIndices[tempRatingIDX]) /
			building.Cell(g.effHeaderIndices[tempRatingIDX])
		cost += building.Cell(g.costHeaderIndices[tempRatingIDX])
	}
	return score
}
func (g *EpcGA) Run(sorter func(candidate1, candidate2 GAState) bool, objective func(gaState GAState) bool) {
	var ph helpers.PrintHelper
	/*
		Create Life! (default-state GAState)
	*/
	var stateRecords []GAStateRecord = make([]GAStateRecord, g.data.Length())
	for i := 0; i < g.data.Length(); i++ {
		stateRecords[i] = CreateGAStateRecord(-1, -1)
	}
	var baseGAState = CreateGAState(stateRecords)

	for i := 0; i < g.maxPopulation; i++ {
		g.population[i] = g.CreateMutation(baseGAState)
	}
	/*
		Do the main process
	*/
	var candidateStates []GAState = make([]GAState, g.maxPopulation*g.childCount+g.maxPopulation)
	//Add existing population to the candidates (immortality exists apparently)
	for i := 0; i < g.maxPopulation; i++ {
		candidateStates[i] = g.population[i]
	}
	var randomInt int
	/*=====================
		Temp, delete log file
	=======================*/

	for roundID := 0; roundID < g.iterations; roundID++ {
		for i := 0; i < len(g.population); i++ {
			for childID := 0; childID < g.childCount; childID++ {
				if rand.Float32() < g.CrossoverRate {
					/*
						Find another GAState to procreate with (bow chaka wow wow!)
					*/
					randomInt = i
					for randomInt == i {
						randomInt = rand.Intn(len(g.population))
					}
					candidateStates[g.maxPopulation+i*g.childCount+childID] = g.Crossover(
						g.population[i],
						g.population[randomInt])
				} else {
					candidateStates[g.maxPopulation+i*g.childCount+childID] = g.CreateMutation(g.population[i])
				}
			}
		}
		/*
			Score the states: Scores are cached in the GAState so you don't need to run it again.
		*/
		for i := 0; i < len(candidateStates); i++ {
			g.scorer.Score(&candidateStates[i]) // Caches results. Doesn't redo every iteration
		}
		/*
			TOURNAMENT TIME!

			Ok, time to literally decimate, or at least n-imate the
			population.

			Simple tournament: There are N competitors but only epcGA.maxPopulation
			spaces in society. Highest epcGA.maxPopulation results get to live
		*/
		sort.Slice(candidateStates, func(i, j int) bool {
			candidate1 := &candidateStates[i]
			candidate2 := &candidateStates[j]

			g.scorer.Score(candidate2) // Caches results. Doesn't redo every iteration
			return g.scorer.Score(candidate1) < g.scorer.Score(candidate2)
		})
		/*
			Sort candidates by objective. This ensures that candidates are sorted by score first,
			then by whether they meet the objective. Doesn't matter for simple objectives like
			score greater than but for thresholds like savedPoints > x.
		*/
		//var meetObjective []GAState = make([]GAState, 0)
		//var doesntMeetObjective []GAState = make([]GAState, 0)
		//for i := 0; i < len(candidateStates); i++ {
		//	if objective(&candidateStates[i]) {
		//		meetObjective = append(meetObjective, candidateStates[i])
		//	} else {
		//		doesntMeetObjective = append(doesntMeetObjective, candidateStates[i])
		//	}
		//}
		//candidateStates = append(meetObjective, doesntMeetObjective...)
		g.population = candidateStates[0 : g.maxPopulation-1]
		/*
			Print and log stuff
		*/
		ph.WriteToFile(g.population[0].ToCSV(), LOG_PATH)
	}
	//ph.P("Scored: " + g.population[0].ToCSV())
	g.population[0].Print()
}
func (g *EpcGA) ScoreGAstate(populationID int) float64 {
	return float64(g.scorer.Score(&g.population[populationID]))
}
func (g *EpcGA) Crossover(state1 GAState, state2 GAState) GAState {

	var states []GAStateRecord = make([]GAStateRecord, g.data.Length())
	for i := 0; i < g.data.Length(); i++ {
		if rand.Float32() < g.CrossoverRate {
			states[i] = CreateGAStateRecord(
				state2.entityStates[i].efficiencyIndex,
				state2.entityStates[i].costIndex)
		} else {
			states[i] = CreateGAStateRecord(
				state1.entityStates[i].efficiencyIndex,
				state1.entityStates[i].costIndex)
		}
	}
	return CreateGAState(states)
}

/*
Create a new mutation
*/
func (g *EpcGA) CreateMutation(baseState GAState) GAState {
	var rnd float32
	var stateIDX int
	var states []GAStateRecord = make([]GAStateRecord, g.data.Length())

	for i := 0; i < g.data.Length(); i++ {
		rnd = rand.Float32()
		if g.Hardness > rnd {

			stateIDX = rand.Intn(len(g.effHeaderIndices)) - 1
			if stateIDX == -1 {
				states[i] = CreateGAStateRecord(
					-1,
					-1)
			} else {
				states[i] = CreateGAStateRecord(
					g.effHeaderIndices[stateIDX],
					g.costHeaderIndices[stateIDX])
			}
		} else {
			states[i] = CreateGAStateRecord(
				baseState.entityStates[i].efficiencyIndex,
				baseState.entityStates[i].costIndex)
		}
	}

	/*
		Send it home
	*/
	return CreateGAState(states)
}
