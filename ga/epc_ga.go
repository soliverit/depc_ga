package ga

import (
	"../csv"
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
	ChildCount        int
	maxPopulation     int
	ForceMaxPoints    bool
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
		ChildCount:		10
		maxPopulation:	30
		hardness:		0.1
	*/

	epcGA.maxPopulation = maxPopulation
	epcGA.population = make([]GAState, epcGA.maxPopulation)
	epcGA.CrossoverRate = 0.15
	epcGA.ChildCount = 10
	epcGA.Hardness = 0.1 //Best 0.1 with crossoverRate of 0.15
	epcGA.ForceMaxPoints = false
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
	/*
		Create Life! (default-state GAState)

		Ok, so it's a multi-objective optimisation problem meaning there's multiple objectives,
		obviously. So, if we start with as-built or random, we might not find any solutions that
		meet the objective. To get round this, we guarantee there's some solutions found by making
		the initial population full of shit but objective-meeting GAState.

		# That said, there needs to be shit options in the start population
	*/
	var stateRecords []GAStateRecord = make([]GAStateRecord, g.data.Length())
	for buildingID := 0; buildingID < g.data.Length(); buildingID++ {

		var bestReduction float32 = 0
		var bestRetrofitID = 0
		for retrofitID := 0; retrofitID < g.data.Building(buildingID).NumberOfRetrofits(); retrofitID++ {
			if g.data.Building(buildingID).Retrofit(retrofitID).Reduction() > bestReduction {
				bestRetrofitID = retrofitID
			}
		}
		if buildingID%2 == 0 {
			stateRecords[buildingID] = CreateGAStateRecord(bestRetrofitID)
		} else {
			stateRecords[buildingID] = CreateGAStateRecord(0)
		}
	}
	var baseGAState = CreateGAState(stateRecords)

	for i := 0; i < g.maxPopulation; i++ {
		g.population[i] = g.CreateMutation(baseGAState)
	}
	/*
		Do the main process
	*/
	var candidateStates []GAState = make([]GAState, g.maxPopulation*g.ChildCount+g.maxPopulation)
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
			for childID := 0; childID < g.ChildCount; childID++ {
				if rand.Float32() < g.CrossoverRate {
					/*
						Find another GAState to procreate with (bow chaka wow wow!)
					*/
					randomInt = i
					for randomInt == i {
						randomInt = rand.Intn(len(g.population))
					}
					candidateStates[g.maxPopulation+i*g.ChildCount+childID] = g.Crossover(
						g.population[i],
						g.population[randomInt])
				} else {
					candidateStates[g.maxPopulation+i*g.ChildCount+childID] = g.CreateMutation(g.population[i])
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
		// NOTE/WARNING: This wouldn't work if we hadn't called ga.scorer.Score() on all candidates a few lines ago. See above
		sort.Slice(candidateStates, func(i, j int) bool {
			return sorter(candidateStates[i], candidateStates[j])
		})
		/*
			Sort candidates by objective. This ensures that candidates are sorted by score first,
			then by whether they meet the objective. Doesn't matter for simple objectives like
			score greater than but for thresholds like savedPoints > x.
		*/
		//var meetObjective []GAState = make([]GAState, 0)
		//for i := 0; i < len(candidateStates); i++ {
		//	if objective(candidateStates[i]) {
		//		meetObjective = append(meetObjective, candidateStates[i])
		//	}
		//}
		//if len(meetObjective) > 0 {
		//	candidateStates = meetObjective
		//}

		g.population = candidateStates[0 : g.maxPopulation-1]
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
			states[i] = CreateGAStateRecord(state2.entityStates[i].optionID)
		} else {
			states[i] = CreateGAStateRecord(state1.entityStates[i].optionID)
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
		building := g.data.Building(i)
		rnd = rand.Float32()
		if g.Hardness > rnd {

			stateIDX = rand.Intn(building.NumberOfRetrofits())
			states[i] = CreateGAStateRecord(stateIDX)
		} else {
			states[i] = CreateGAStateRecord(baseState.entityStates[i].optionID)
		}
	}

	/*
		Send it home
	*/
	return CreateGAState(states)
}
