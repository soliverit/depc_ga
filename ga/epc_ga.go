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
	Create a new EpcGA, an extension of BaseGA for residential EPC data
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
func (ga *EpcGA) Best() float32 {
	var building *csv.Building
	/*
		Get header positions
	*/
	var cost float32 = 0.0
	var score float32 = 0.0
	var tempScore float32 = 0.0
	var tempRatingIDX int = 0
	var testScore float32
	for i := 0; i < ga.data.Length(); i++ {
		building = ga.data.Building(i)
		for j := 0; j < len(ga.effHeaders); j++ {
			testScore = building.Cell(ga.effHeaderIndices[j]) / building.Cell(ga.costHeaderIndices[j])
			if testScore > tempScore {
				tempScore = testScore
				tempRatingIDX = j
			}
		}
		score += building.Cell(ga.costHeaderIndices[tempRatingIDX]) /
			building.Cell(ga.effHeaderIndices[tempRatingIDX])
		cost += building.Cell(ga.costHeaderIndices[tempRatingIDX])
	}
	return score
}
func (ga *EpcGA) Run() {
	var ph helpers.PrintHelper
	/*
		Create Life! (default-state GAState)
	*/
	var stateRecords []GAStateRecord = make([]GAStateRecord, ga.data.Length())
	for i := 0; i < ga.data.Length(); i++ {
		stateRecords[i] = CreateGAStateRecord(-1, -1)
	}
	var baseGAState = CreateGAState(stateRecords)

	for i := 0; i < ga.maxPopulation; i++ {
		ga.population[i] = ga.CreateMutation(baseGAState)
	}
	/*
		Do the main process
	*/
	var candidateStates []GAState = make([]GAState, ga.maxPopulation*ga.childCount+ga.maxPopulation)
	//Add existing population to the candidates (immortality exists apparently)
	for i := 0; i < ga.maxPopulation; i++ {
		candidateStates[i] = ga.population[i]
	}
	var randomInt int
	/*=====================
		Temp, delete log file
	=======================*/

	for roundID := 0; roundID < ga.iterations; roundID++ {
		for i := 0; i < len(ga.population); i++ {
			for childID := 0; childID < ga.childCount; childID++ {
				if rand.Float32() < ga.CrossoverRate {
					/*
						Find another GAState to procreate with (bow chaka wow wow!)
					*/
					randomInt = i
					for randomInt == i {
						randomInt = rand.Intn(len(ga.population))
					}
					candidateStates[ga.maxPopulation+i*ga.childCount+childID] = ga.Crossover(
						ga.population[i],
						ga.population[randomInt])
				} else {
					candidateStates[ga.maxPopulation+i*ga.childCount+childID] = ga.CreateMutation(ga.population[i])
				}
			}
		}
		/*
			TOURNAMENT TIME!

			Ok, time to literally decimate, or at least n-imate the
			population.

			Simple tournament: There are N competitors but only EpcGA.maxPopulation
			spaces in society. Highest EpcGA.maxPopulation results get to live
		*/
		sort.Slice(candidateStates, func(i, j int) bool {
			return ga.scorer.Score(&candidateStates[i]) < ga.scorer.Score(&candidateStates[j])
		})
		ga.population = candidateStates[0 : ga.maxPopulation-1]
		/*
			Print and log stuff
		*/

		ph.WriteToFile(ga.population[0].ToCSV(), LOG_PATH)
	}
	ph.P("Scored: " + ga.population[0].ToCSV())
}
func (ga *EpcGA) ScoreGAstate(populationID int) float64 {
	return float64(ga.scorer.Score(&ga.population[populationID]))
}
func (ga *EpcGA) Crossover(state1 GAState, state2 GAState) GAState {

	var states []GAStateRecord = make([]GAStateRecord, ga.data.Length())
	for i := 0; i < ga.data.Length(); i++ {
		if rand.Float32() < ga.CrossoverRate {
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
func (ga *EpcGA) CreateMutation(baseState GAState) GAState {
	var rnd float32
	var stateIDX int
	var states []GAStateRecord = make([]GAStateRecord, ga.data.Length())

	for i := 0; i < ga.data.Length(); i++ {
		rnd = rand.Float32()
		if ga.Hardness > rnd {

			stateIDX = rand.Intn(len(ga.effHeaderIndices)) - 1
			if stateIDX == -1 {
				states[i] = CreateGAStateRecord(
					-1,
					-1)
			} else {
				states[i] = CreateGAStateRecord(
					ga.effHeaderIndices[stateIDX],
					ga.costHeaderIndices[stateIDX])
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
