package main

import (
	"./csv"
	"./ga"
	"./helpers"
	"./optimiser"
	bayesopt "go-bayesopt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

const PATH string = "C:\\workspaces\\__shared_data__\\depc\\"
const CERTS string = "certificates.csv"
const OUTPUT string = "indices.csv"
const RETROFITS string = "retrofits.csv"
const LOCATION string = PATH + "domestic-E06000002-Middlesbrough\\"
const GA_PATH string = LOCATION + RETROFITS
const GA_TARGETS string = LOCATION + "targets.csv"

var RETROFIT_LABELS = []string{"envelopes",
	"hotwater_envelopes",
	"hotwater_roof_envelopes",
	"hotwater_roof_windows",
	"hotwater_windows",
	"hotwater",
	"roof_envelopes",
	"roof_hotwater",
	"roof_windows",
	"roof",
	"windows_envelopes",
	"windows_hotwater_envelopes",
	"windows_hotwater_roof_envelopes",
	"windows_roof_envelopes",
	"windows"}
var RETROFIT_TARGET_LABELS = []string{
	"hotwater",
	"hotwater_envelopes",
	"roof_envelopes",
	"windows_envelopes",
	"envelopes",
	"roof_hotwater",
	"hotwater_windows",
	"hotwater",
	"roof_windows",
	"roof",
	"windows"}

/*
dhw, dhw-env, dhw-win, dhw-roof
roof, roof-env, roof-win,
*/
const OPTIMISE bool = false

func main() {

	// Clear console
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
	// Set rand seed
	rand.Seed(1)
	/*
		Do helpers
	*/
	var ph *helpers.PrintHelper = helpers.CreatePrintHelper(true)
	ph.Line("====== Starting process using data from: " + GA_PATH + " ======")

	/*
		Parse main data
	*/
	var data csv.BuildingReader = csv.ParseBuildingCSV(GA_PATH, false)
	ph.Line("Data loaded - " + strconv.Itoa(data.Length()) + " records")

	/*
		Add original values (epc index, co2, energy) and Subtract original ratings from all values
	*/
	var targets csv.BuildingReader = csv.ParseBuildingCSV(GA_TARGETS, false)
	var targetLength int = targets.Length()
	var targetCo2IDX int = targets.EasyReader.ColumnNameToIndex("co2Emissions")
	var targetsEfficiencyIDX int = targets.EasyReader.ColumnNameToIndex("energyEfficiency")
	var subtractColumns []int = make([]int, len(RETROFIT_TARGET_LABELS))
	for i := 0; i < len(subtractColumns); i++ {
		subtractColumns[i] = data.EasyReader.ColumnNameToIndex(RETROFIT_LABELS[i] + "-Eff")
	}
	/*
		Link target data
	*/
	for i := 0; i < targetLength; i++ {
		if data.Building(i).Cell(0) != -9999 {
			targetRow := targets.Row(i)
			//data.Building(i).Subtract(targetRow.Cell(targetsEfficiencyIDX), subtractColumns)
			data.Building(i).SetOriginalProperties(targetRow.Cell(targetsEfficiencyIDX), targetRow.Cell(targetCo2IDX))
		}
	}
	/*
		Remove corrupt Buildings
	*/
	data.RemoveCorrupt()
	ph.Line("Data cleansed - " + strconv.Itoa(data.Length()) + " records")
	ph.Line("WARNING!!!: Sampling 0.5 of records. OOPS!")
	data = *data.Sample(0.5)
	data.PrepareRetrofits(RETROFIT_LABELS)
	if OPTIMISE {
		helpers.PrintRow([]string{"Score", "Points", "Cost", "P/C", "C/P"}, 15)

		var bayesOpt optimiser.GAOptimiser = optimiser.CreateBayesOptimiser()
		bayesOpt.AddParam("Hardness", 0.05, 0.15)
		bayesOpt.AddParam("CrossoverRate", 0.05, 0.2)
		bayesOpt.Run(func(params map[bayesopt.Param]float64) float64 {
			/*
				Reset random number

				TODO: Is this serial? How does random state consistent over threads
			*/
			rand.Seed(1)
			var epcGA *ga.EpcGA = ga.CreateEpcGA(&data, 100, 40, RETROFIT_TARGET_LABELS)

			for param, value := range params {
				switch param.GetName() {
				case "Hardness":
					epcGA.Hardness = float32(value)
				case "CrossoverRate":
					epcGA.CrossoverRate = float32(value)
				}
			}
			epcGA.Run(func(candidate1, candidate2 ga.GAState) bool {
				return candidate1.Score() < candidate2.Score()
			}, func(candidate ga.GAState) bool {
				return candidate.Points() < 100000
			})
			/*
				fmt.Println("V:" + strconv.FormatFloat(epcGA.ScoreGAstate(0),'f',2,64))
			*/
			return epcGA.ScoreGAstate(0)
		})

	} else {
		/*====
		From Bayes-Opt

		Hardness:		0.110466
		CrossoverRate:	0.191076
		*/
		/*
			Create newEpcGA
		*/
		var epcGA *ga.EpcGA = ga.CreateEpcGA(&data, 1000, 100, RETROFIT_TARGET_LABELS)
		epcGA.Hardness = 0.110466
		epcGA.CrossoverRate = 0.191076
		epcGA.ChildCount = 20
		ph.Line("EPC-GA instantiated")
		/*
			Get minimum index point improvement from
		*/
		var baseEPC float32
		for buildingID := 0; buildingID < data.Length(); buildingID++ {
			baseEPC += data.Building(buildingID).EPCIndex()
		}

		/*
			This is the GA entry point for the case study
		*/
		data.WriteToFile("c:/workspaces/__shared_data__/depc_ga_buildings_test.csv")
		var targetEPC float32 = baseEPC * 0.25
		println("Target EPC: " + strconv.FormatFloat(float64(targetEPC), 'f', 0, 32))
		epcGA.Run(func(candidate1, candidate2 ga.GAState) bool {
			return candidate1.Score() < candidate2.Score()
		}, func(candidate ga.GAState) bool {
			return candidate.Points() < 120000
		})
	}
}
func doWeatherNN() {
	/*
		data path
	*/
	var pathPattern string = "c:\\repos\\__shared_data__\\ml_anon\\*\\model_epc.inp"
	paths, _ := filepath.Glob(pathPattern)
	for i := 0; i < len(paths); i++ {

	}

}
