package main

import (
	"./csv"
	"./ga"
	"./helpers"
	"./optimiser"
	bayesopt "go-bayesopt"
	"math/rand"
	"path/filepath"
	"strconv"
	"strings"
)

const PATH string = "C:\\workspaces\\__shared_data__\\depc\\"
const CERTS string = "certificates.csv"
const OUTPUT string = "indices.csv"
const RETROFITS string = "retrofits.csv"
const LOCATION string = PATH + "domestic-E06000002-Middlesbrough\\"
const GA_PATH string = LOCATION + RETROFITS
const GA_TARGETS string = LOCATION + "targets.csv"

var RETROFIT_LABELS = []string{"envelopes_hotwater_roof_windows",
	"envelopes_hotwater_roof",
	"envelopes_hotwater_windows",
	"envelopes_hotwater",
	"envelopes_roof_windows",
	"envelopes_roof",
	"envelopes_windows",
	"envelopes",
	"hotwater_roof_windows",
	"hotwater_roof",
	"hotwater_windows",
	"hotwater",
	"roof_windows",
	"roof",
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
		Subtract original ratings from all values
	*/
	var targets csv.BuildingReader = csv.ParseBuildingCSV(GA_TARGETS, false)
	var targetLength int = targets.Length()
	var targetsColumnIDX int = targets.EasyReader.ColumnNameToIndex("energyEfficiency")
	var subtractColumns []int = make([]int, len(RETROFIT_TARGET_LABELS))
	for i := 0; i < len(subtractColumns); i++ {
		subtractColumns[i] = data.EasyReader.ColumnNameToIndex(RETROFIT_TARGET_LABELS[i] + "-Eff")
		/*
			!!! Patch boiler cost !!!
		*/
		if strings.Contains(RETROFIT_TARGET_LABELS[i], "hotwater") {
			for buildingID := 0; buildingID < data.Length(); buildingID++ {
				data.Building(buildingID).Subtract(-1200,
					[]int{data.EasyReader.ColumnNameToIndex(RETROFIT_TARGET_LABELS[i] + "-Cost")})
			}
		}
	}
	for i := 0; i < targetLength; i++ {
		if data.Building(i).Cell(0) != -9999 {
			data.Building(i).Subtract(targets.Row(i).Cell(targetsColumnIDX), subtractColumns)
		}
	}
	/*
		Remove corrupt Buildings
	*/
	data.RemoveCorrupt()
	ph.Line("Data cleansed - " + strconv.Itoa(data.Length()) + " records")
	data = *data.Sample(0.5)
	if OPTIMISE {
		var bayesOpt optimiser.GAOptimiser = optimiser.CreateBayesOptimiser()
		bayesOpt.AddParam("Hardness", 0.05, 0.15)
		bayesOpt.AddParam("CrossoverRate", 0.05, 0.2)
		bayesOpt.Run(func(params map[bayesopt.Param]float64) float64 {
			/*
				Reset random number

				TODO: Is this serial? How does random state consistent over threads
			*/
			rand.Seed(1)
			var epcGA *ga.EpcGA = ga.CreateEpcGA(&data, 99, 40, RETROFIT_TARGET_LABELS)
			for param, value := range params {
				switch param.GetName() {
				case "Hardness":
					epcGA.Hardness = float32(value)
				case "CrossoverRate":
					epcGA.CrossoverRate = float32(value)
				}
			}
			epcGA.Run(func(candidate1, candidate2 ga.GAState) bool {
				return true
			}, func(candidate ga.GAState) bool { return true })
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
		var epcGA *ga.EpcGA = ga.CreateEpcGA(&data, 100, 40, RETROFIT_TARGET_LABELS)
		epcGA.Hardness = 0.110466
		epcGA.CrossoverRate = 0.191076
		ph.Line("EPC-GA instantiated")
		/*
			This is the GA entry point for the case study
		*/
		epcGA.Run(func(candidate1, candidate2 ga.GAState) bool {
			return candidate1.Score() > 0
		}, func(candidate ga.GAState) bool { return true })
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
